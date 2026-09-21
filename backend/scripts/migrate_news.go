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
}
