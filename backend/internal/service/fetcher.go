package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
				query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4)`
				_, err = s.DB.Exec(context.Background(), query, news.Category, news.Title, news.Summary, "BMKG")
				if err != nil {
					log.Printf("[Zero-Admin] Gagal menyimpan berita BMKG ke DB: %v\n", err)
				} else {
					log.Println("[Zero-Admin] Berita BMKG berhasil disimpan ke Database.")
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

	client := &http.Client{Timeout: 10 * time.Second}
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
			log.Println("[Zero-Admin] Berhasil menarik data JSON Vulkanik dari NASA EONET. Meringkas dengan AI...")

			// Ambil 5 event pertama saja agar AI tidak kelebihan context limit
			var shortData interface{} = data
			if events, ok := data["events"].([]interface{}); ok {
				if len(events) > 5 {
					shortData = map[string]interface{}{"events": events[:5]}
				}
			}
			rawDataBytes, _ := json.Marshal(shortData)

			news, err := ai.GenerateNewsSummary(context.Background(), string(rawDataBytes), "NASA EONET (Earth Observatory Natural Event Tracker)")

			if err != nil {
				log.Printf("[Zero-Admin] Gagal meringkas berita NASA: %v\n", err)
			} else if s.DB != nil {
				// Simpan hasil ringkasan ke database
				query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4)`
				_, err = s.DB.Exec(context.Background(), query, news.Category, news.Title, news.Summary, "NASA EONET")
				if err != nil {
					log.Printf("[Zero-Admin] Gagal menyimpan berita NASA ke DB: %v\n", err)
				} else {
					log.Println("[Zero-Admin] Berita NASA berhasil disimpan ke Database.")
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
