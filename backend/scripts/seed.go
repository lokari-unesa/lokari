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

// seed.go — pengisian & perawatan data potensi_bencana (membutuhkan
// COHERE_API_KEY di .env seperti seed sebelumnya).
//
//	go run scripts/seed.go           # seed dari Posko.txt (lokasi + vektor AI)
//	go run scripts/seed.go backfill  # lengkapi baris yang embedding-nya masih NULL (idempotent)
//	go run scripts/seed.go dedupe    # bersihkan baris duplikat potensi_bencana (per nama_objek)

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

	cmd := "seed"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "seed":
		seedPotensi(ctx, pool)
	case "backfill":
		backfillEmbedding(ctx, pool)
	case "dedupe":
		dedupePotensi(ctx, pool)
	default:
		log.Fatalf("Argumen tidak dikenal: %q — pilih: seed | backfill | dedupe\n", cmd)
	}
}

// seedPotensi membaca scripts/Posko.txt lalu menyimpan lokasi + vektor AI-nya.
func seedPotensi(ctx context.Context, pool *pgxpool.Pool) {
	filePath := "scripts/Posko.txt"
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Gagal membaca file Posko.txt: %v", err)
	}

	// Perbaiki beberapa array JSON yang tergabung (mis. "]\n[").
	content := string(data)
	content = strings.ReplaceAll(content, "]\r\n[", ",")
	content = strings.ReplaceAll(content, "]\n[", ",")
	content = strings.ReplaceAll(content, "][", ",")

	var nodes []ClaudeNode
	if err := json.Unmarshal([]byte(content), &nodes); err != nil {
		log.Fatalf("Gagal parsing JSON: %v", err)
	}

	log.Printf("Berhasil membaca %d lokasi dari Posko.txt", len(nodes))
	log.Println("Mulai membangkitkan vektor AI untuk setiap lokasi...")

	count := 0
	for _, node := range nodes {
		deskripsiFull := fmt.Sprintf("%s Berada di %s. Fasilitas ini masuk kategori %s. Mampu menampung sekitar %d pengungsi.",
			node.Deskripsi, node.AlamatDusun, node.Kategori, node.KapasitasOrang)

		// Cegah duplikasi
		var exists bool
		err := pool.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM potensi_bencana WHERE nama_objek = $1)",
			node.NamaObjek).Scan(&exists)
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

		_, err = pool.Exec(ctx, `
			INSERT INTO potensi_bencana (nama_objek, kategori, deskripsi, alamat_dusun, kapasitas_orang, geometri, embedding)
			VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8)`,
			node.NamaObjek, node.Kategori, node.Deskripsi, node.AlamatDusun,
			node.KapasitasOrang, node.Lon, node.Lat, pgvector.NewVector(emb))

		if err != nil {
			log.Printf("Gagal insert database: %s: %v\n", node.NamaObjek, err)
		} else {
			log.Printf("Berhasil menyimpan: %s", node.NamaObjek)
			count++
		}

		time.Sleep(500 * time.Millisecond) // hindari rate limit tier gratis
	}

	log.Printf("TUNTAS! %d lokasi dari Claude berhasil ditanam di sistem Vektor LOKARI!", count)
}

// backfillEmbedding melengkapi baris potensi_bencana yang embedding-nya NULL
// agar terjangkau semantic search. Idempotent (baris berembedding dilewati).
func backfillEmbedding(ctx context.Context, pool *pgxpool.Pool) {
	type lokasi struct {
		ID             string
		NamaObjek      string
		Kategori       string
		Deskripsi      string
		AlamatDusun    string
		KapasitasOrang int
	}

	rows, err := pool.Query(ctx, `
		SELECT id_potensi, nama_objek, kategori, deskripsi, alamat_dusun,
		       COALESCE(kapasitas_orang, 0)
		FROM potensi_bencana
		WHERE embedding IS NULL`)
	if err != nil {
		log.Fatalf("Gagal mengambil baris tanpa embedding: %v\n", err)
	}

	var targets []lokasi
	for rows.Next() {
		var l lokasi
		if err := rows.Scan(&l.ID, &l.NamaObjek, &l.Kategori, &l.Deskripsi, &l.AlamatDusun, &l.KapasitasOrang); err != nil {
			log.Fatalf("Gagal membaca baris: %v\n", err)
		}
		targets = append(targets, l)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		log.Fatalf("Error saat iterasi baris: %v\n", err)
	}

	if len(targets) == 0 {
		log.Println("Tidak ada baris tanpa embedding. Selesai.")
		return
	}
	log.Printf("%d baris tanpa embedding ditemukan. Mulai backfill...\n", len(targets))

	updated := 0
	for _, l := range targets {
		// Teks embedding dibuat sama seperti seedPotensi agar vektornya konsisten.
		deskripsiFull := fmt.Sprintf("%s Berada di %s. Fasilitas ini masuk kategori %s. Mampu menampung sekitar %d pengungsi.",
			l.Deskripsi, l.AlamatDusun, l.Kategori, l.KapasitasOrang)

		emb, err := ai.GenerateEmbedding(l.NamaObjek + ". " + deskripsiFull)
		if err != nil {
			log.Printf("Gagal generate embedding untuk %s: %v. Lanjut...\n", l.NamaObjek, err)
			time.Sleep(1 * time.Second)
			continue
		}

		_, err = pool.Exec(ctx,
			`UPDATE potensi_bencana SET embedding = $2, updated_at = CURRENT_TIMESTAMP WHERE id_potensi = $1`,
			l.ID, pgvector.NewVector(emb))
		if err != nil {
			log.Printf("Gagal update %s: %v\n", l.NamaObjek, err)
			continue
		}

		updated++
		log.Printf("[%d/%d] Backfill: %s", updated, len(targets), l.NamaObjek)
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("TUNTAS! %d dari %d baris berhasil dilengkapi embedding.", updated, len(targets))
}

// dedupePotensi menghapus baris ganda potensi_bencana: disisakan baris paling
// awal (created_at terendah, lalu id_potensi terkecil) per nama_objek.
func dedupePotensi(ctx context.Context, pool *pgxpool.Pool) {
	log.Println("Menghapus data ganda (duplikat) di tabel potensi_bencana...")

	tag, err := pool.Exec(ctx, `
		DELETE FROM potensi_bencana
		WHERE id_potensi NOT IN (
			SELECT id_potensi FROM (
				SELECT id_potensi,
					ROW_NUMBER() OVER (PARTITION BY nama_objek ORDER BY created_at, id_potensi) AS rn
				FROM potensi_bencana
			) t WHERE rn = 1
		);`)
	if err != nil {
		log.Fatalf("Gagal menghapus duplikat: %v", err)
	}

	log.Printf("Selesai! %d baris data ganda berhasil dibersihkan.", tag.RowsAffected())
}