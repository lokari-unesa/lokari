//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"

	"os"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	_ = godotenv.Load(".env", "../.env")
	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to lokari_db: %v\n", err)
	}
	defer pool.Close()

	log.Println("Menghapus data ganda (duplicates) di tabel potensi_bencana...")

	// Sisakan baris PALING AWAL (created_at terendah, lalu id_potensi terkecil
	// sebagai tie-breaker) untuk tiap nama_objek — bukan id_potensi::text acak.
	query := `
		DELETE FROM potensi_bencana
		WHERE id_potensi NOT IN (
			SELECT id_potensi FROM (
				SELECT id_potensi,
					ROW_NUMBER() OVER (PARTITION BY nama_objek ORDER BY created_at, id_potensi) AS rn
				FROM potensi_bencana
			) t WHERE rn = 1
		);
	`
	tag, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Gagal menghapus duplikat: %v", err)
	}
	
	log.Printf("Selesai! %d baris data ganda berhasil dibersihkan.", tag.RowsAffected())
}
