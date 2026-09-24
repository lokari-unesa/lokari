//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// migrate.go — satu pintu untuk migrasi & reset skema database LOKARI.
//
//	go run scripts/migrate.go           # terapkan SEMUA migrasi (idempotent, aman diulang)
//	go run scripts/migrate.go reset     # DESTRUKTIF: hapus semua tabel lalu bangun ulang
//
// Membutuhkan DATABASE_URL di .env (sudah tersedia di environment container).

func main() {
	_ = godotenv.Load(".env", "../.env")
	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}

	if len(os.Args) > 1 && os.Args[1] == "reset" {
		resetDB(lokariDB)
		return
	}
	migrateAll(lokariDB)
}

// migrateAll menerapkan seluruh skema (idempotent, aman diulang berkali-kali).
func migrateAll(lokariDB string) {
	ctx := context.Background()
	ensureDatabase(ctx, lokariDB)

	pool, err := pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to lokari_db: %v\n", err)
	}
	defer pool.Close()

	migrateSchema(ctx, pool)
	log.Println("Migrasi selesai: skema database LOKARI lengkap (idempotent).")
}

func ensureDatabase(ctx context.Context, lokariDB string) {
	defaultDB := strings.Replace(lokariDB, "lokari_db", "postgres", 1)
	pool, err := pgxpool.New(ctx, defaultDB)
	if err != nil {
		log.Printf("Tidak dapat konek ke database postgres (barangkali memang belum ada): %v\n", err)
		return
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, "CREATE DATABASE lokari_db"); err != nil {
		log.Printf("Database lokari_db mungkin sudah ada: %v\n", err)
	} else {
		log.Println("Database lokari_db dibuat.")
	}
}

// migrateSchema membuat ekstensi + seluruh tabel + index-nya. Semua perintah
// memakai IF NOT EXISTS sehingga aman dijalankan ulang kapan pun.
func migrateSchema(ctx context.Context, pool *pgxpool.Pool) {
	log.Println("Mengaktifkan ekstensi (postgis, uuid-ossp, vector)...")
	for _, ext := range []string{
		"CREATE EXTENSION IF NOT EXISTS postgis;",
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";",
		"CREATE EXTENSION IF NOT EXISTS vector;",
	} {
		if _, err := pool.Exec(ctx, ext); err != nil {
			log.Fatalf("Gagal mengaktifkan ekstensi: %v\n", err)
		}
	}

	log.Println("Membuat tabel inti (kategori_layer, potensi_bencana, log_update)...")
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS kategori_layer (
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
    embedding vector(1024), -- vektor dari Cohere embed-multilingual-v3.0 (1024 dimensi)
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
);`)
	if err != nil {
		log.Fatalf("Error membuat tabel inti: %v\n", err)
	}

	log.Println("Membuat index spasial & vektor (GiST, HNSW)...")
	if _, err := pool.Exec(ctx,
		`CREATE INDEX IF NOT EXISTS idx_potensi_geom ON potensi_bencana USING GiST (geometri);
		 CREATE INDEX IF NOT EXISTS idx_potensi_embed ON potensi_bencana USING hnsw (embedding vector_cosine_ops);`); err != nil {
		log.Fatalf("Gagal membuat index: %v\n", err)
	}

	log.Println("Membuat tabel berita (kabar_kelud) + index unik (sumber, judul)...")
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS kabar_kelud (
    id_kabar UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kategori VARCHAR(50) NOT NULL,
    judul VARCHAR(255) NOT NULL,
    ringkasan TEXT NOT NULL,
    sumber VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		log.Fatalf("Error membuat tabel kabar_kelud: %v\n", err)
	}

	// Bersihkan duplikat lama (sisakan baris paling awal per sumber+judul) agar
	// index unik berikutnya tidak gagal karena baris ganda yang sudah ada.
	if _, err := pool.Exec(ctx, `
		DELETE FROM kabar_kelud
		WHERE id_kabar NOT IN (
			SELECT id_kabar FROM (
				SELECT id_kabar,
					ROW_NUMBER() OVER (PARTITION BY sumber, judul ORDER BY created_at, id_kabar) AS rn
				FROM kabar_kelud
			) t WHERE rn = 1
		);`); err != nil {
		log.Fatalf("Error membersihkan berita duplikat: %v\n", err)
	}
	if _, err := pool.Exec(ctx,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_kabar_kelud_sumber_judul ON kabar_kelud (sumber, judul)`); err != nil {
		log.Fatalf("Error membuat index unik kabar_kelud: %v\n", err)
	}

	log.Println("Membuat tabel push (push_subscriptions)...")
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS push_subscriptions (
    id SERIAL PRIMARY KEY,
    endpoint TEXT UNIQUE NOT NULL,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		log.Fatalf("Error membuat tabel push_subscriptions: %v\n", err)
	}

	log.Println("Membuat tabel state monitor (monitor_state)...")
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS monitor_state (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		log.Fatalf("Error membuat tabel monitor_state: %v\n", err)
	}
}

// resetDB menghapus SEMUA tabel lalu membangun ulang skema dari nol,
// dengan konfirmasi manual dan seed ulang opsional.
func resetDB(lokariDB string) {
	if !confirm("Ini akan MENGHAPUS SEMUA DATA di database. Lanjut? [y/N] ") {
		log.Println("Dibatalkan. Tidak ada perubahan.")
		return
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	log.Println("Menghapus semua tabel...")
	_, err = pool.Exec(ctx, `
		DROP TABLE IF EXISTS kabar_kelud CASCADE;
		DROP TABLE IF EXISTS monitor_state CASCADE;
		DROP TABLE IF EXISTS push_subscriptions CASCADE;
		DROP TABLE IF EXISTS potensi_bencana CASCADE;
		DROP TABLE IF EXISTS kategori_layer CASCADE;
		DROP TABLE IF EXISTS log_update CASCADE;
	`)
	if err != nil {
		log.Fatalf("Gagal menghapus tabel: %v\n", err)
	}

	migrateAll(lokariDB)

	if confirm("Isi ulang data dengan seed (membutuhkan COHERE_API_KEY)? [y/N] ") {
		runStep("go", "run", "scripts/seed.go")
	} else {
		log.Println("Seed dilewati. Jalankan manual: docker compose exec backend go run scripts/seed.go")
	}

	log.Println("RESET SUKSES! Skema dibangun ulang dari nol.")
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func runStep(name string, args ...string) {
	log.Printf("Menjalankan: %s %s ...", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Gagal menjalankan %s %s: %v\n", name, strings.Join(args, " "), err)
	}
}