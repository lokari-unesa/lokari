//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	lokariDB := "postgresql://postgres:root@127.0.0.1:5433/lokari_db?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to lokari_db: %v\n", err)
	}
	defer pool.Close()

	log.Println("Menghapus data ganda (duplicates) di tabel potensi_bencana...")
	
	// Query to keep only the earliest inserted row for each nama_objek
	query := `
		DELETE FROM potensi_bencana
		WHERE id_potensi NOT IN (
			SELECT MIN(id_potensi::text)::uuid
			FROM potensi_bencana
			GROUP BY nama_objek
		);
	`
	tag, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Gagal menghapus duplikat: %v", err)
	}
	
	log.Printf("Selesai! %d baris data ganda berhasil dibersihkan.", tag.RowsAffected())
}
