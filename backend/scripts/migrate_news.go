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

	log.Println("Creating table kabar_kelud...")

	schema := `
CREATE TABLE IF NOT EXISTS kabar_kelud (
    id_kabar UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kategori VARCHAR(50) NOT NULL,
    judul VARCHAR(255) NOT NULL,
    ringkasan TEXT NOT NULL,
    sumber VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	_, err = pool.Exec(ctx, schema)
	if err != nil {
		log.Fatalf("Error creating table: %v\n", err)
	}
	log.Println("Table kabar_kelud created successfully.")

	// Bersihkan duplikat lama (sisakan baris paling awal per kombinasi sumber+judul).
	// Wajib dijalankan sebelum unique index dibuat agar index tidak gagal karena
	// baris ganda yang sudah terlanjur ada.
	dupCleanup := `
		DELETE FROM kabar_kelud
		WHERE id_kabar NOT IN (
			SELECT id_kabar FROM (
				SELECT id_kabar,
					ROW_NUMBER() OVER (PARTITION BY sumber, judul ORDER BY created_at, id_kabar) AS rn
				FROM kabar_kelud
			) t WHERE rn = 1
		);
	`
	tag, err := pool.Exec(ctx, dupCleanup)
	if err != nil {
		log.Fatalf("Error cleaning duplicate news: %v\n", err)
	}
	if tag.RowsAffected() > 0 {
		log.Printf("%d baris berita duplikat dibersihkan.\n", tag.RowsAffected())
	}

	// Cegah duplikasi berikutnya: satu judul per sumber hanya boleh muncul sekali.
	_, err = pool.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_kabar_kelud_sumber_judul ON kabar_kelud (sumber, judul)`)
	if err != nil {
		log.Fatalf("Error creating unique index: %v\n", err)
	}
	log.Println("Unique index kabar_kelud (sumber, judul) terpasang.")
}
