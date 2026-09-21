//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lokari/backend/internal/ai"
	"github.com/pgvector/pgvector-go"
)

type ClaudeNode struct {
	NamaObjek      string  `json:"nama_objek"`
	Kategori       string  `json:"kategori"`
	Deskripsi      string  `json:"deskripsi"`
	AlamatDusun    string  `json:"alamat_dusun"`
	KapasitasOrang int     `json:"kapasitas_orang"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
}

func main() {
	_ = godotenv.Load(".env", "../.env")

	if os.Getenv("COHERE_API_KEY") == "" {
		log.Fatal("ERROR: COHERE_API_KEY belum diatur di .env")
	}

	lokariDB := "postgresql://postgres:root@127.0.0.1:5433/lokari_db?sslmode=disable"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, lokariDB)
	if err != nil {
		log.Fatalf("Unable to connect to lokari_db: %v\n", err)
	}
	defer pool.Close()

	// Read Posko.txt
	filePath := "../frontend/static/assets/Posko.txt"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Gagal membaca file Posko.txt: %v", err)
	}

	// Fix multiple JSON arrays concatenated together (e.g. "]\n[")
	content := string(data)
	content = strings.ReplaceAll(content, "]\r\n[", ",")
	content = strings.ReplaceAll(content, "]\n[", ",")
	content = strings.ReplaceAll(content, "][", ",")

	var nodes []ClaudeNode
	if err := json.Unmarshal([]byte(content), &nodes); err != nil {
		log.Fatalf("Gagal parsing JSON: %v", err)
	}

	log.Printf("Berhasil membaca %d lokasi dari Posko.txt", len(nodes))
	log.Println("Mulai membangkitkan Kecerdasan Buatan (Vektor) untuk setiap lokasi...")

	count := 0

	for _, node := range nodes {
		// Konteks untuk Search
		deskripsiFull := fmt.Sprintf("%s Berada di %s. Fasilitas ini masuk kategori %s. Mampu menampung sekitar %d pengungsi.", 
			node.Deskripsi, node.AlamatDusun, node.Kategori, node.KapasitasOrang)

		// Cegah duplikasi
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM potensi_bencana WHERE nama_objek = $1)", node.NamaObjek).Scan(&exists)
		if err == nil && exists {
			log.Printf("⏩ Melewati %s (Sudah ada di database)", node.NamaObjek)
			continue
		}

		log.Printf("[%d/%d] Memproses Vektor AI: %s ...", count+1, len(nodes), node.NamaObjek)

		emb, err := ai.GenerateEmbedding(node.NamaObjek + ". " + deskripsiFull)
		if err != nil {
			log.Printf("Gagal Vektor AI untuk %s: %v. Lanjut...", node.NamaObjek, err)
			time.Sleep(1 * time.Second)
			continue
		}

		vec := pgvector.NewVector(emb)

		sql := `
			INSERT INTO potensi_bencana (nama_objek, kategori, deskripsi, alamat_dusun, kapasitas_orang, geometri, embedding)
			VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8)
		`
		_, err = pool.Exec(ctx, sql, node.NamaObjek, node.Kategori, node.Deskripsi, node.AlamatDusun, node.KapasitasOrang, node.Lon, node.Lat, vec)

		if err != nil {
			log.Printf("Gagal insert database: %s: %v\n", node.NamaObjek, err)
		} else {
			log.Printf("Berhasil menyimpan: %s", node.NamaObjek)
			count++
		}

		// Delay to avoid hitting rate limits on free tier
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("TUNTAS! %d lokasi dari Claude berhasil ditanam di sistem Vektor LOKARI!", count)
}
