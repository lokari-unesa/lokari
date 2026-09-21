package worker

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/service"
	"github.com/robfig/cron/v3"
)

func StartCronJobs(db *pgxpool.Pool) *cron.Cron {
	c := cron.New()
	fetcher := service.NewFetcherService(db)

	log.Println("[Zero-Admin] Menginisialisasi jadwal Cron Worker...")

	// Jalankan setiap 6 jam (Sesuai kesepakatan untuk menghemat kuota API)
	// Format Cron: Menit Jam Tanggal Bulan Hari
	_, err := c.AddFunc("0 */6 * * *", func() {
		fetcher.FetchBMKGData()
		fetcher.FetchNASAData()
	})

	if err != nil {
		log.Fatalf("Gagal menjadwalkan Cron: %v", err)
	}

	// Jalankan sekali saat server baru menyala (agar tidak menunggu 1 jam penuh untuk data pertama)
	go fetcher.FetchBMKGData()
	go fetcher.FetchNASAData()

	c.Start()
	log.Println("[Zero-Admin] Cron Worker berhasil dijalankan di latar belakang.")
	
	return c
}
