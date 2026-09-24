//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lokari/backend/internal/ai"
	"github.com/pgvector/pgvector-go"
)

// backfill_embedding.go — melengkapi baris potensi_bencana yang embedding-nya
// masih NULL (mis. hasil input manual atau seed sebelum fitur vektor ada),
// agar baris tersebut terjangkau oleh semantic search.
//
// Idempotent: baris yang sudah punya embedding dibiarkan. Membutuhkan
// COHERE_API_KEY seperti seed_claude.go.
func main() {
	_ = godotenv.Load(".env", "../.env")

	if os.Getenv("COHERE_API_KEY") == "" {
		log.Fatal("ERROR: COHERE_API_KEY belum diatur di .env")
	}
	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	type lokasi struct {
		ID            string
		NamaObjek     string
		Kategori      string
		Deskripsi     string
		AlamatDusun   string
		KapasitasOrang int
	}

	// Baris tanpa embedding: embeddig disimpan sebagai vector, jadi NULL.
	rows, err := pool.Query(ctx, `
		SELECT id_potensi, nama_objek, kategori, deskripsi, alamat_dusun,
		       COALESCE(kapasitas_orang, 0)
		FROM potensi_bencana
		WHERE embedding IS NULL
	`)
	if err != nil {
		log.Fatalf("Gagal mengambil baris tanpa embedding: %v\n", err)
	}

	var targets []lokasi
	for rows.Next() {
		var l lokasi
		if err := rows.Scan(&l.ID, &l.NamaObjek, &l.Kategori, &l.Deskripsi, &l.AlamatDusun, &l.KapasitasOrang); err != nil {
			log.Fatalf("Gagal membaca baris: %v\n", err)
		}
		targets = append(targets, l)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		log.Fatalf("Error saat iterasi baris: %v\n", err)
	}

	if len(targets) == 0 {
		log.Println("Tidak ada baris tanpa embedding. Selesai.")
		return
	}
	log.Printf("%d baris tanpa embedding ditemukan. Mulai backfill...\n", len(targets))

	updated := 0
	for _, l := range targets {
		// Bangun teks embedding dengan format yang sama seperti seed_claude.go
		// agar vektornya konsisten.
		deskripsiFull := fmt.Sprintf("%s Berada di %s. Fasilitas ini masuk kategori %s. Mampu menampung sekitar %d pengungsi.",
			l.Deskripsi, l.AlamatDusun, l.Kategori, l.KapasitasOrang)

		emb, err := ai.GenerateEmbedding(l.NamaObjek + ". " + deskripsiFull)
		if err != nil {
			log.Printf("Gagal generate embedding untuk %s: %v. Lanjut...\n", l.NamaObjek, err)
			time.Sleep(1 * time.Second)
			continue
		}

		vec := pgvector.NewVector(emb)
		_, err = pool.Exec(ctx,
			`UPDATE potensi_bencana SET embedding = $2, updated_at = CURRENT_TIMESTAMP WHERE id_potensi = $1`,
			l.ID, vec)
		if err != nil {
			log.Printf("Gagal update %s: %v\n", l.NamaObjek, err)
			continue
		}

		updated++
		log.Printf("[%d/%d] Backfill: %s", updated, len(targets), l.NamaObjek)
		time.Sleep(500 * time.Millisecond) // hindari rate limit tier gratis
	}

	log.Printf("TUNTAS! %d dari %d baris berhasil dilengkapi embedding.", updated, len(targets))
}