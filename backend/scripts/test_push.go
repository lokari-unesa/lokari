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

	"github.com/SherClockHolmes/webpush-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env", "../.env")

	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}

	vapidPublic := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivate := os.Getenv("VAPID_PRIVATE_KEY")
	if vapidPublic == "" || vapidPrivate == "" {
		log.Fatal("ERROR: VAPID keys belum diatur di .env. Generate: go run scripts/vapid.go")
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

	payload, _ := json.Marshal(map[string]string{
		"title": "LOKARI: Notifikasi Uji Coba",
		"body":  "Selamat, notifikasi LOKARI berfungsi dengan baik.",
		"url":   "/news",
	})

	sent, deleted := 0, 0
	for _, sub := range subs {
		resp, err := webpush.SendNotification(payload, sub, &webpush.Options{
			Subscriber:      "mailto:admin@lokari.my.id",
			VAPIDPublicKey:  vapidPublic,
			VAPIDPrivateKey: vapidPrivate,
			TTL:             60, // uji coba — cukup bertahan 1 menit
		})
		switch {
		case err != nil:
			log.Printf("  gagal (network): %s → %v\n", sub.Endpoint, err)
		case resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound:
			if _, delErr := pool.Exec(ctx,
				"DELETE FROM push_subscriptions WHERE endpoint = $1", sub.Endpoint); delErr == nil {
				deleted++
				log.Printf("  subscription mati dihapus: %s\n", sub.Endpoint)
			}
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			sent++
			log.Printf("  terkirim: %s\n", sub.Endpoint)
		default:
			log.Printf("  gagal (HTTP %d): %s\n", resp.StatusCode, sub.Endpoint)
		}
	}

	fmt.Printf("\nHasil: %d terkirim, %d subscription mati dibersihkan, %d total terdaftar.\n", sent, deleted, len(subs))
}
