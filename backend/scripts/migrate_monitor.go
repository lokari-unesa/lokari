//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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

	schema := `
CREATE TABLE IF NOT EXISTS monitor_state (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	_, err = pool.Exec(ctx, schema)
	if err != nil {
		log.Fatalf("Error creating monitor_state table: %v\n", err)
	}
	log.Println("Database schema for monitor_state migrated successfully.")
}