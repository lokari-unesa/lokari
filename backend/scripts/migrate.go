//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"

	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env", "../.env")
	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}
	defaultDB := strings.Replace(lokariDB, "lokari_db", "postgres", 1)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, defaultDB)
	if err != nil {
		log.Fatalf("Unable to connect to default database: %v\n", err)
	}

	// Try to create DB
	_, err = pool.Exec(ctx, "CREATE DATABASE lokari_db")
	if err != nil {
		log.Printf("Database lokari_db might already exist or error: %v\n", err)
	} else {
		log.Println("Database lokari_db created successfully.")
	}
	pool.Close()

	// Now connect to lokari_db
	pool, err = pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to lokari_db: %v\n", err)
	}
	defer pool.Close()

	// Enable postgis, uuid-ossp, and vector extensions
	log.Println("Enabling required extensions (postgis, uuid-ossp, vector)...")
	extensions := []string{
		"CREATE EXTENSION IF NOT EXISTS postgis;",
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";",
		"CREATE EXTENSION IF NOT EXISTS vector;",
	}
	for _, ext := range extensions {
		if _, err := pool.Exec(ctx, ext); err != nil {
			log.Fatalf("Gagal mengaktifkan ekstensi: %v\n", err)
		}
	}

	log.Println("Creating tables (potensi_bencana, kategori_layer, log_update)...")

	// Create tables
	schema := `

CREATE TABLE IF NOT EXISTS kategori_layer (
    id_kategori SERIAL PRIMARY KEY,
    nama_kategori VARCHAR(100) NOT NULL,
    ikon_marker VARCHAR(255)
);


CREATE TABLE IF NOT EXISTS potensi_bencana (
    id_potensi UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_objek VARCHAR(255) NOT NULL,
    kategori VARCHAR(100) NOT NULL,
    tingkat_risiko VARCHAR(50),
    deskripsi TEXT,
    alamat_dusun VARCHAR(255),
    kapasitas_orang INTEGER,
    geometri GEOMETRY,
    embedding vector(1024),
    foto_lokasi VARCHAR(255),
    kontak_darurat VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS log_update (
    id_log UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sumber_api VARCHAR(100) NOT NULL,
    status_tarik VARCHAR(50) NOT NULL,
    waktu_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
	_, err = pool.Exec(ctx, schema)
	if err != nil {
		log.Fatalf("Error creating tables: %v\n", err)
	}

	log.Println("Creating indexes (GiST geometric, HNSW embedding)...")

	// Index spasial PostGIS untuk pencarian berbasis jarak/dekat
	idxGeom := `CREATE INDEX IF NOT EXISTS idx_potensi_geom ON potensi_bencana USING GiST (geometri);`
	if _, err := pool.Exec(ctx, idxGeom); err != nil {
		log.Fatalf("Gagal membuat index GiST: %v\n", err)
	}

	// Index vektor HNSW untuk semantic search (cosine)
	idxEmbed := `CREATE INDEX IF NOT EXISTS idx_potensi_embed ON potensi_bencana USING hnsw (embedding vector_cosine_ops);`
	if _, err := pool.Exec(ctx, idxEmbed); err != nil {
		log.Fatalf("Gagal membuat index HNSW: %v\n", err)
	}

	log.Println("Database schema migrated successfully.")
}
