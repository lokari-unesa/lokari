package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/ai"
)

type FetcherService struct {
	DB *pgxpool.Pool
}

func NewFetcherService(db *pgxpool.Pool) *FetcherService {
	return &FetcherService{DB: db}
}

// FetchBMKGData pulls JSON data from BMKG Open Data API
func (s *FetcherService) FetchBMKGData() {
	log.Println("[Zero-Admin] Memulai penarikan data gempabumi BMKG...")
	start := time.Now()

	// BMKG Official JSON API for Earthquakes
	bmkgURL := "https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json"
	status := "Sukses"

	req, err := http.NewRequest("GET", bmkgURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[Zero-Admin] BMKG Fetch Error: %v\n", err)
		status = "Gagal"
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("[Zero-Admin] BMKG responded with status: %d\n", resp.StatusCode)
		status = "Gagal"
	} else {
		defer resp.Body.Close()

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("[Zero-Admin] BMKG JSON Parse Error: %v\n", err)
			status = "Gagal"
		} else {
			log.Println("[Zero-Admin] Berhasil menarik data JSON Gempa dari BMKG API. Meringkas dengan AI...")

			// Ubah data map ke JSON string untuk dilempar ke AI
			rawDataBytes, _ := json.Marshal(data)
			news, err := ai.GenerateNewsSummary(context.Background(), string(rawDataBytes), "BMKG (Badan Meteorologi, Klimatologi, dan Geofisika)")

			if err != nil {
				log.Printf("[Zero-Admin] Gagal meringkas berita BMKG: %v\n", err)
			} else if s.DB != nil {
				// Simpan hasil ringkasan ke database
				query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4) ON CONFLICT (sumber, judul) DO NOTHING`
				_, err = s.DB.Exec(context.Background(), query, news.Category, news.Title, news.Summary, "BMKG")
				if err != nil {
					log.Printf("[Zero-Admin] Gagal menyimpan berita BMKG ke DB: %v\n", err)
				} else {
					log.Println("[Zero-Admin] Berita BMKG berhasil disimpan ke Database.")
					if news.Category == "warning" || news.Category == "volcano" || news.Category == "evac" {
						go s.broadcastPush(news.Title, news.Summary, news.Category)
					}
				}
			}
		}
	}

	// Log to database
	s.logUpdate("BMKG", status)
	log.Printf("[Zero-Admin] Tarik BMKG Selesai: %s (%v)\n", status, time.Since(start))
}

// FetchNASAData pulls volcanic event data from NASA EONET (Earth Observatory Natural Event Tracker)
func (s *FetcherService) FetchNASAData() {
	log.Println("[Zero-Admin] Memulai penarikan data aktivitas vulkanik global dari API NASA EONET...")
	start := time.Now()

	// NASA EONET API for active volcanoes
	nasaURL := "https://eonet.gsfc.nasa.gov/api/v3/events?category=volcanoes&status=open"
	status := "Sukses"

	// EONET sering respons lambat (>10 detik), beri kelonggaran 60 detik.
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(nasaURL)
	if err != nil {
		log.Printf("[Zero-Admin] NASA EONET Fetch Error: %v\n", err)
		status = "Gagal"
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("[Zero-Admin] NASA EONET responded with status: %d\n", resp.StatusCode)
		status = "Gagal"
	} else {
		defer resp.Body.Close()
		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("[Zero-Admin] NASA EONET JSON Parse Error: %v\n", err)
			status = "Gagal"
		} else {
			log.Println("[Zero-Admin] Berhasil menarik data JSON Vulkanik dari NASA EONET.")

			// Hanya proses event yang berada dalam radius 150 km dari Gunung
			// Kelud — event global yang jauh tidak relevan untuk warga Kediri.
			shortData, err := eonetNearKeludData(data)
			if err != nil {
				log.Printf("[Zero-Admin] NASA EONET: %v\n", err)
			} else {
				rawDataBytes, _ := json.Marshal(shortData)

				news, err := ai.GenerateNewsSummary(context.Background(), string(rawDataBytes), "NASA EONET (Earth Observatory Natural Event Tracker)")
				if err != nil {
					log.Printf("[Zero-Admin] Gagal meringkas berita NASA: %v\n", err)
				} else if s.DB != nil {
					// Simpan hasil ringkasan ke database
					query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4) ON CONFLICT (sumber, judul) DO NOTHING`
					_, err = s.DB.Exec(context.Background(), query, news.Category, news.Title, news.Summary, "NASA EONET")
					if err != nil {
						log.Printf("[Zero-Admin] Gagal menyimpan berita NASA ke DB: %v\n", err)
					} else {
						log.Println("[Zero-Admin] Berita NASA berhasil disimpan ke Database.")
						if news.Category == "warning" || news.Category == "volcano" || news.Category == "evac" {
							go s.broadcastPush(news.Title, news.Summary, news.Category)
						}
					}
				}
			}
		}
	}

	s.logUpdate("NASA_EONET", status)
	log.Printf("[Zero-Admin] Tarik NASA EONET Selesai: %s (%v)\n", status, time.Since(start))
}

// logUpdate mencatat riwayat pembaruan ke tabel log_update
func (s *FetcherService) logUpdate(sumberAPI string, status string) {
	if s.DB == nil {
		return
	}

	query := `INSERT INTO log_update (sumber_api, status_tarik) VALUES ($1, $2)`
	_, err := s.DB.Exec(context.Background(), query, sumberAPI, status)
	if err != nil {
		log.Printf("[Zero-Admin] Gagal mencatat log_update: %v\n", err)
	}
}

// broadcastPush menembakkan notifikasi Web Push ke seluruh token warga yang terdaftar
func (s *FetcherService) broadcastPush(title, summary, category string) {
	if s.DB == nil {
		return
	}
	rows, err := s.DB.Query(context.Background(), "SELECT endpoint, p256dh, auth FROM push_subscriptions")
	if err != nil {
		log.Println("[Push] Gagal mengambil list warga:", err)
		return
	}
	defer rows.Close()

	vapidPublic := os.Getenv("VAPID_PUBLIC_KEY")
	vapidPrivate := os.Getenv("VAPID_PRIVATE_KEY")
	if vapidPublic == "" || vapidPrivate == "" {
		log.Println("[Push] VAPID Keys belum diatur di .env. Notifikasi dibatalkan.")
		return
	}

	payloadMap := map[string]string{
		"title": "LOKARI " + strings.ToUpper(category) + ": " + title,
		"body":  summary,
		"url":   "/news",
	}
	b, _ := json.Marshal(payloadMap)

	count := 0
	deleted := 0
	for rows.Next() {
		var ep, p256, auth string
		if err := rows.Scan(&ep, &p256, &auth); err != nil {
			continue
		}

		sub := &webpush.Subscription{
			Endpoint: ep,
			Keys: webpush.Keys{P256dh: p256, Auth: auth},
		}

		resp, err := webpush.SendNotification(b, sub, &webpush.Options{
			Subscriber:      "mailto:admin@lokari.my.id",
			VAPIDPublicKey:  vapidPublic,
			VAPIDPrivateKey: vapidPrivate,
			Urgency:         webpush.UrgencyHigh, // Info penting — jangan ditunda push service
			TTL:             3600,                // Aktif selama 1 jam
		})
		switch {
		case err != nil:
			// Gagal level transpor (network/timeout) — bukan salah subscription, biarkan
			continue
		case resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound:
			// Subscription sudah tidak valid (dihapus browser / kedaluwarsa) — buang dari DB
			if _, delErr := s.DB.Exec(context.Background(),
				"DELETE FROM push_subscriptions WHERE endpoint = $1", ep); delErr == nil {
				deleted++
			}
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			count++
		default:
			log.Printf("[Push] Gagal kirim ke %s (HTTP %d)\n", ep, resp.StatusCode)
		}
	}
	log.Printf("[Push] Selesai: %d notifikasi terkirim, %d subscription mati dibersihkan\n", count, deleted)
}
