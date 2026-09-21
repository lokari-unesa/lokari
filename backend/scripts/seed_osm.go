//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lokari/backend/internal/ai"
	"github.com/pgvector/pgvector-go"
)

// Struct untuk mem-parsing JSON dari Overpass API
type OSMResponse struct {
	Elements []OSMNode `json:"elements"`
}

type OSMNode struct {
	Type   string            `json:"type"`
	Id     int64             `json:"id"`
	Lat    float64           `json:"lat"`
	Lon    float64           `json:"lon"`
	Center *struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"center"`
	Tags map[string]string `json:"tags"`
}

func main() {
	// Load .env tanpa memunculkan log error yang mengganggu jika salah satu file tidak ada
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

	// 1. Koordinat Desa Jarak, Kediri: Lat -7.952, Lon 112.235 (Radius 8km untuk mencegah masuk wilayah Malang/Blitar)
	// Kita ganti 'node' menjadi 'nwr' (node, way, relation) untuk menangkap bangunan (area/polygon) seperti MI
	// Dan tambahkan 'out center;' agar polygon memiliki titik tengah (center)
	query := `[out:json];(nwr(around:8000,-7.952,112.235)[amenity~"school|hospital|clinic|community_centre|place_of_worship"];nwr(around:8000,-7.952,112.235)[building~"school|public|civic"];);out center;`

	apiURL := "http://overpass-api.de/api/interpreter"
	log.Println("Sedang menarik data spasial nyata dari Satelit OpenStreetMap (Radius 8km Desa Jarak, Kediri)...")
	
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader("data="+url.QueryEscape(query)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "LOKARI-Dev/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Gagal menarik data dari OSM: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var osmData OSMResponse
	if err := json.Unmarshal(bodyBytes, &osmData); err != nil {
		log.Fatalf("Gagal mem-parsing data OSM: %v", err)
	}

	nodes := osmData.Elements
	if len(nodes) == 0 {
		log.Fatal("Tidak ada fasilitas publik yang ditemukan di radius 15km.")
	}
	log.Printf("Ditemukan %d bangunan fasilitas publik nyata dari satelit!", len(nodes))
	log.Println("Mulai membangkitkan Kecerdasan Buatan (Vektor) untuk setiap lokasi. Proses ini memakan waktu beberapa menit, harap sabar...")

	// Kita batasi maksimal 30 lokasi saja agar proses Cohere API tidak terlalu lama untuk demonstrasi
	limit := 30
	count := 0

	for _, node := range nodes {
		if count >= limit {
			break
		}

		nama := node.Tags["name"]
		if nama == "" {
			continue // Skip bangunan tak bernama
		}

		amenity := node.Tags["amenity"]
		kategori := "Titik Kumpul Umum"
		kapasitas := 50

		switch {
		case amenity == "school" || node.Tags["building"] == "school":
			kategori = "Posko Pendidikan (Sekolah / Madrasah)"
			kapasitas = 200
		case amenity == "hospital" || amenity == "clinic":
			kategori = "Posko Kesehatan"
			kapasitas = 100
		case amenity == "community_centre" || node.Tags["building"] == "public" || node.Tags["building"] == "civic":
			kategori = "Balai Desa / Gedung Serbaguna"
			kapasitas = 300
		case amenity == "place_of_worship":
			kategori = "Tempat Ibadah (Masjid/Gereja)"
			kapasitas = 150
		}

		// Konteks Desa Jarak untuk Search
		deskripsi := fmt.Sprintf("Bangunan %s ini adalah %s nyata di lapangan. Dapat dijadikan tempat aman atau titik kumpul evakuasi bagi warga yang tinggal di kawasan bahaya erupsi Gunung Kelud (khususnya wilayah Desa Jarak, Dusun Jarak Lor, Jarak Kidul, Sagi, Kalasan Barat, Kalasan Timur dan sekitarnya).", nama, kategori)

		// Cegah duplikasi data di database
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM potensi_bencana WHERE nama_objek = $1)", nama).Scan(&exists)
		if err == nil && exists {
			log.Printf("⏩ Melewati %s (Sudah ada di database)", nama)
			continue
		}

		log.Printf("[%d/%d] Memproses Vektor AI: %s ...", count+1, limit, nama)

		// Generate embedding via Cohere
		emb, err := ai.GenerateEmbedding(nama + ". " + deskripsi)
		if err != nil {
			log.Printf("Gagal Vektor AI untuk %s: %v. Lanjut ke lokasi berikutnya...\n", nama, err)
			time.Sleep(1 * time.Second) // Hindari rate limit parah
			continue
		}

		vec := pgvector.NewVector(emb)

		lat := node.Lat
		lon := node.Lon
		if (node.Type == "way" || node.Type == "relation") && node.Center != nil {
			lat = node.Center.Lat
			lon = node.Center.Lon
		}

		// Insert dengan PostGIS geometri
		sql := `
			INSERT INTO potensi_bencana (nama_objek, kategori, deskripsi, kapasitas_orang, geometri, embedding)
			VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($5, $6), 4326), $7)
		`
		_, err = pool.Exec(ctx, sql, nama, kategori, deskripsi, kapasitas, lon, lat, vec)

		if err != nil {
			log.Printf("Gagal insert database: %s: %v\n", nama, err)
		} else {
			log.Printf("Berhasil menyimpan: %s", nama)
			count++
		}

		// Delay kecil untuk menghormati Rate Limit Cohere Free Tier
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("TUNTAS! %d lokasi nyata (Real Live) berhasil tertanam di sistem Vektor LOKARI!", count)
}
