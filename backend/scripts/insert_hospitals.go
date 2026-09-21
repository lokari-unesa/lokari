//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
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

	// RSUD Gambiran
	sql1 := `
		INSERT INTO potensi_bencana (nama_objek, kategori, deskripsi, alamat_dusun, kapasitas_orang, geometri)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326))
		ON CONFLICT DO NOTHING
	`
	_, err = pool.Exec(ctx, sql1, 
		"RSUD Gambiran Kota Kediri", 
		"Fasilitas Kesehatan", 
		"Rumah Sakit Umum Daerah kelas B di pusat Kota Kediri. Fasilitas medis paling lengkap di wilayah Kediri Raya, sangat aman dari ancaman lahar karena berada puluhan kilometer di luar zona merah.", 
		"Jl. Kapten Tendean No.16, Kota Kediri", 
		1500, 
		112.0298, 
		-7.8488)
	if err != nil {
		log.Printf("Gagal insert RSUD Gambiran: %v\n", err)
	} else {
		fmt.Println("Berhasil insert RSUD Gambiran")
	}

	// RSKK Pare
	_, err = pool.Exec(ctx, sql1, 
		"RSUD Kabupaten Kediri (RSKK Pare)", 
		"Fasilitas Kesehatan", 
		"Rumah Sakit Umum Daerah milik Pemerintah Kabupaten Kediri di Pare. Pusat rujukan utama bencana alam Kabupaten Kediri, berada di zona yang sangat aman dari jangkauan awan panas dan lahar Kelud.", 
		"Jl. Pahlawan Kusuma Bangsa No.1, Pare", 
		1000, 
		112.1813, 
		-7.7617)
	if err != nil {
		log.Printf("Gagal insert RSKK Pare: %v\n", err)
	} else {
		fmt.Println("Berhasil insert RSKK Pare")
	}
}
