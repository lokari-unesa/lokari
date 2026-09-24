//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lokari/backend/internal/service"
)

// ops.go — utilitas operasional harian LOKARI.
//
//	go run scripts/ops.go hot|cold|all  # jalankan job penarikan data secara manual (tanpa menunggu cron)
//	go run scripts/ops.go push          # kirim DEMO semua jenis notif (data dummy) ke seluruh subscriber
//	go run scripts/ops.go vapid         # generate pasangan kunci VAPID baru
//
// Membutuhkan DATABASE_URL di .env (sudah tersedia di environment container);
// mode push juga butuh VAPID_PUBLIC_KEY & VAPID_PRIVATE_KEY.

func main() {
	_ = godotenv.Load(".env", "../.env")

	if len(os.Args) < 2 {
		log.Fatal("Argumen wajib: hot | cold | all | push | vapid")
	}

	switch os.Args[1] {
	case "hot", "cold", "all":
		triggerLoops(os.Args[1])
	case "push":
		testPush()
	case "vapid":
		generateVAPID()
	default:
		log.Fatalf("Argumen tidak dikenal: %q — pilih: hot | cold | all | push | vapid\n", os.Args[1])
	}
}

// triggerLoops menjalankan job penarikan data secara manual — berguna saat
// server baru dinyalakan (belum ada data) atau untuk tes, tanpa menunggu
// jadwal cron (hot = 30 detik, cold = 6 jam).
func triggerLoops(mode string) {
	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}

	db, err := pgxpool.New(context.Background(), lokariDB)
	if err != nil {
		log.Fatalf("ERROR: gagal konek database: %v\n", err)
	}
	defer db.Close()

	f := service.NewFetcherService(db)

	switch mode {
	case "hot":
		log.Println("Trigger hot loop: gempa BMKG + status Gunung Kelud...")
		f.FetchBMKGFeltQuakes()
		f.FetchKeludStatus()
	case "cold":
		log.Println("Trigger cold loop: laporan harian MAGMA + event NASA EONET...")
		f.FetchMagmaLaporan()
		f.FetchNASAData()
	case "all":
		log.Println("Trigger hot loop: gempa BMKG + status Gunung Kelud...")
		f.FetchBMKGFeltQuakes()
		f.FetchKeludStatus()
		log.Println("Trigger cold loop: laporan harian MAGMA + event NASA EONET...")
		f.FetchMagmaLaporan()
		f.FetchNASAData()
	}
	log.Println("Selesai.")
}

// testPush mengirim notifikasi uji coba ke seluruh subscriber & membersihkan
// subscription yang sudah tidak valid (HTTP 410/404).
func testPush() {
	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}

	vapidPublic := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivate := os.Getenv("VAPID_PRIVATE_KEY")
	if vapidPublic == "" || vapidPrivate == "" {
		log.Fatal("ERROR: VAPID keys belum diatur di .env. Generate: go run scripts/ops.go vapid")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to lokari_db: %v\n", err)
	}
	defer pool.Close()

	rows, err := pool.Query(ctx, "SELECT endpoint, p256dh, auth FROM push_subscriptions")
	if err != nil {
		log.Fatalf("Gagal membaca push_subscriptions: %v\n", err)
	}
	defer rows.Close()

	var subs []*webpush.Subscription
	for rows.Next() {
		var ep, p256, auth string
		if err := rows.Scan(&ep, &p256, &auth); err != nil {
			continue
		}
		subs = append(subs, &webpush.Subscription{
			Endpoint: ep,
			Keys:     webpush.Keys{P256dh: p256, Auth: auth},
		})
	}
	if rows.Err() != nil {
		log.Fatalf("Gagal membaca baris subscription: %v\n", rows.Err())
	}

	if len(subs) == 0 {
		log.Println("Belum ada subscription terdaftar. Buka frontend (HTTPS/localhost) lalu izinkan notifikasi browser.")
		return
	}

	// Data dummy yang meniru SEMUA jenis notifikasi yang mampu dikirim sistem:
	//   1. Perubahan status Kelud NAIK  → kategori volcano (MAGMA)
	//   2. Perubahan status Kelud TURUN → kategori volcano (MAGMA)
	//   3. Gempa BMKG yang terasa       → kategori warning (BMKG)
	//   4. Event vulkanik NASA EONET    → kategori volcano (NASA EONET)
	// Format judul & body sengaja dibuat sama persis dengan broadcastPush asli
	// ("LOKARI <KATEGORI>: ...") supaya tampilan popup sama seperti kondisi nyata.
	type demoNotif struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		URL   string `json:"url"`
	}

	demos := []demoNotif{
		{
			Title: "LOKARI VOLCANO: Status Gunung Kelud: Level II (Waspada)",
			Body:  "Tingkat aktivitas Gunung Kelud berubah: Level I (Normal) → Level II (Waspada). Jauhi aliran sungai dan pantau terus informasi resmi PVMBG MAGMA serta arahan BPBD setempat.",
			URL:   "/news",
		},
		{
			Title: "LOKARI VOLCANO: Status Gunung Kelud: Level I (Normal)",
			Body:  "Tingkat aktivitas Gunung Kelud berubah: Level II (Waspada) → Level I (Normal). Warga dapat beraktivitas seperti biasa, namun tetap waspada.",
			URL:   "/news",
		},
		{
			Title: "LOKARI WARNING: GEMPA M 4,2 TERASA — Kab. Kediri",
			Body:  "Gempa berkekuatan M 4,2 terjadi 85 km arah barat daya Kota Kediri, dirasakan MMI III di Desa Jarak. Tetap tenang, jauhi bangunan yang retak, dan ikuti arahan petugas.",
			URL:   "/news",
		},
		{
			Title: "LOKARI VOLCANO: Awan Panas Gunung Semeru Terpantau",
			Body:  "NASA EONET mencatat aktivitas vulkanik di Gunung Semeru (± 77 km dari Kediri). Hindari area rawan dan ikuti arahan BPBD setempat.",
			URL:   "/news",
		},
	}

	sent, deleted := 0, 0
	for i, demo := range demos {
		payload, _ := json.Marshal(demo)

		for _, sub := range subs {
			resp, err := webpush.SendNotification(payload, sub, &webpush.Options{
				Subscriber:      "mailto:admin@lokari.my.id",
				VAPIDPublicKey:  vapidPublic,
				VAPIDPrivateKey: vapidPrivate,
				Urgency:         webpush.UrgencyHigh, // Info penting — jangan ditunda push service
				TTL:             3600,                // sama seperti notif asli (aktif 1 jam)
			})
			switch {
			case err != nil:
				log.Printf("  [%d/4] gagal (network): %s → %v\n", i+1, sub.Endpoint, err)
			case resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound:
				if _, delErr := pool.Exec(ctx,
					"DELETE FROM push_subscriptions WHERE endpoint = $1", sub.Endpoint); delErr == nil {
					deleted++
					log.Printf("  subscription mati dihapus: %s\n", sub.Endpoint)
				}
			case resp.StatusCode >= 200 && resp.StatusCode < 300:
				sent++
			default:
				log.Printf("  [%d/4] gagal (HTTP %d): %s\n", i+1, resp.StatusCode, sub.Endpoint)
			}
		}
		log.Printf("[%d/4] Terkirim: %s\n", i+1, demo.Title)

		// Jeda supaya popup muncul satu per satu, tidak menumpuk.
		if i < len(demos)-1 {
			log.Println("Jeda 6 detik sebelum notif berikutnya...")
			time.Sleep(6 * time.Second)
		}
	}

	fmt.Printf("\nHasil: %d kiriman terkirim, %d subscription mati dibersihkan, %d total terdaftar.\n", sent, deleted, len(subs))
}

// generateVAPID menampilkan pasangan kunci VAPID baru untuk .env.
func generateVAPID() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== ADD THESE TO YOUR .ENV FILE ===")
	fmt.Println("VAPID_PUBLIC_KEY=" + publicKey)
	fmt.Println("VAPID_PRIVATE_KEY=" + privateKey)
	fmt.Println("===================================")
}