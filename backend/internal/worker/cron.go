package worker

import (
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/service"
	"github.com/robfig/cron/v3"
)

// StartCronJobs menjadwalkan dua siklus penarikan data:
//
//	Hot loop  — tiap 30 detik: gempa BMKG terbaru & status Gunung Kelud
//	            (ancaman yang butuh respon cepat / peringatan dini).
//	Cold loop — tiap 6 jam: laporan harian MAGMA & event NASA EONET
//	            (informasi pelengkap, hemat kuota API).
func StartCronJobs(db *pgxpool.Pool) *cron.Cron {
	c := cron.New()
	fetcher := service.NewFetcherService(db)

	log.Println("[Zero-Admin] Menginisialisasi jadwal Cron Worker...")

	// Mutex mencegah dua siklus hot loop bertabrakan bila satu fetch lebih
	// lambat dari interval 30 detik (tick berikutnya dilewati, tidak menumpuk).
	var hotMu sync.Mutex
	hotLoop := func() {
		if !hotMu.TryLock() {
			return
		}
		defer hotMu.Unlock()
		fetcher.FetchBMKGFeltQuakes()
		fetcher.FetchKeludStatus()
	}

	var coldMu sync.Mutex
	coldLoop := func() {
		if !coldMu.TryLock() {
			return
		}
		defer coldMu.Unlock()
		fetcher.FetchMagmaLaporan()
		fetcher.FetchNASAData()
	}

	if _, err := c.AddFunc("@every 30s", hotLoop); err != nil {
		log.Fatalf("Gagal menjadwalkan hot loop: %v", err)
	}

	if _, err := c.AddFunc("0 */6 * * *", coldLoop); err != nil {
		log.Fatalf("Gagal menjadwalkan cold loop: %v", err)
	}

	c.Start()
	log.Println("[Zero-Admin] Cron Worker berjalan: hot loop 30 detik, cold loop 6 jam.")
	return c
}
