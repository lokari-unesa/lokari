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

// reset.go — operasi DESTRUKTIF. Menghapus SEMUA tabel lalu membangun ulang
// skema (migrate.go, migrate_news.go) dan opsional mengisi seed ulang.
//
// Berbeda dengan migrate.go yang aman (non-destruktif, idempotent), skrip ini
// WAJIB mendapat konfirmasi manual sebelum menghapus data apa pun.
func main() {
	_ = godotenv.Load(".env", "../.env")

	lokariDB := os.Getenv("DATABASE_URL")
	if lokariDB == "" {
		log.Fatal("ERROR: DATABASE_URL is required in .env")
	}

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

	log.Println("Menghapus tabel yang ada (kabar_kelud, potensi_bencana, kategori_layer, log_update)...")
	_, err = pool.Exec(ctx, `
		DROP TABLE IF EXISTS kabar_kelud CASCADE;
		DROP TABLE IF EXISTS potensi_bencana CASCADE;
		DROP TABLE IF EXISTS kategori_layer CASCADE;
		DROP TABLE IF EXISTS log_update CASCADE;
	`)
	if err != nil {
		log.Fatalf("Gagal menghapus tabel: %v\n", err)
	}

	// Bangun ulang skema lewat skrip migrasi yang sama (satu sumber kebenaran).
	runStep("go", "run", "scripts/migrate.go")
	runStep("go", "run", "scripts/migrate_news.go")

	if confirm("Isi ulang data dengan seed (membutuhkan COHERE_API_KEY)? [y/N] ") {
		runStep("go", "run", "scripts/seed_claude.go")
	} else {
		log.Println("Seed dilewati. Jalankan manual: docker compose exec backend go run scripts/seed_claude.go")
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