package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/ai"
	"github.com/pgvector/pgvector-go"
)

// AIHandler handles all AI-related routes
type AIHandler struct {
	DB *pgxpool.Pool
}

func NewAIHandler(db *pgxpool.Pool) *AIHandler {
	return &AIHandler{DB: db}
}

// GetAlert fetches NASA EONET data and uses DeepSeek to summarize it
func (h *AIHandler) GetAlert(c *fiber.Ctx) error {
	// 1. Fetch real-time NASA EONET data — dengan timeout & cek status code
	nasaURL := "https://eonet.gsfc.nasa.gov/api/v3/events?category=volcanoes&status=open&limit=5"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(nasaURL)
	if err != nil {
		log.Printf("[handler:alert] NASA fetch error: %v", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Gagal menarik data satelit NASA"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		log.Printf("[handler:alert] NASA responded with status %d: %s", resp.StatusCode, string(body))
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Data satelit NASA sedang tidak tersedia"})
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		log.Printf("[handler:alert] NASA body read error: %v", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Gagal membaca data satelit NASA"})
	}
	rawData := string(bodyBytes)

	// 2. Pass the raw JSON to our DeepSeek NLP Summarizer
	summary, err := ai.SummarizeAlert(context.Background(), rawData)
	if err != nil {
		if strings.Contains(err.Error(), "KUOTA_OPENROUTER_HABIS") {
			return c.JSON(fiber.Map{
				"sumber": "error.quota.alert.src",
				"pesan":  "error.quota.alert.msg",
			})
		}
		log.Printf("[handler:alert] summarize error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI gagal merangkum pesan"})
	}

	return c.JSON(fiber.Map{
		"sumber": "NASA EONET (Diproses oleh AI DeepSeek)",
		"pesan":  summary,
	})
}

// SearchRequest is the expected JSON body for Semantic Search
type SearchRequest struct {
	Query string   `json:"query"`
	Lat   *float64 `json:"lat"`
	Lng   *float64 `json:"lng"`
}

// SearchSemantic performs vector search on potensi_bencana using pgvector
func (h *AIHandler) SearchSemantic(c *fiber.Ctx) error {
	if h.DB == nil {
		log.Println("[handler:search] DB pool nil — cek DATABASE_URL di .env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak terhubung"})
	}

	var req SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input JSON tidak valid"})
	}

	// Validasi query server-side — jangan hanya andalkan blocklist client
	req.Query = strings.TrimSpace(req.Query)
	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query tidak boleh kosong"})
	}
	if len(req.Query) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query terlalu panjang (maksimal 500 karakter)"})
	}
	for _, r := range req.Query {
		if r < 32 && r != '\t' {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query mengandung karakter yang tidak valid"})
		}
	}

	// 1. Convert user's sentence to a vector embedding using HuggingFace
	embedding, err := ai.GenerateEmbedding(req.Query)
	if err != nil {
		if strings.Contains(err.Error(), "KUOTA_COHERE_HABIS") {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "KUOTA_HABIS", "message": "error.quota.search.msg"})
		}
		log.Printf("[handler:search] embedding error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat vektor embedding"})
	}

	// 2. Tentukan Strategi Pencarian: RETRIEVE & RERANK
	// Backend hanya bertugas melakukan RETRIEVE: Mengambil Top 15 hasil yang secara makna (semantik) paling relevan.
	// Frontend nanti yang akan melakukan RERANK: Mengurutkan ulang Top 15 ini berdasarkan jarak fisik (jika ada kata 'dekat').

	sqlQuery := `
		SELECT 
			id_potensi, 
			nama_objek, 
			kategori, 
			deskripsi,
			kapasitas_orang,
			ST_AsGeoJSON(geometri) as geojson,
			(embedding <=> $1) as distance 
		FROM potensi_bencana 
		ORDER BY embedding <=> $1 
		LIMIT 15;
	`

	vec := pgvector.NewVector(embedding)
	rows, err := h.DB.Query(context.Background(), sqlQuery, vec)
	if err != nil {
		return serverError(c, "search.query", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	var skipped int
	for rows.Next() {
		var id, nama, kategori, geojson string
		var deskripsi *string // NULL di DB → dikirim null, bukan di-skip
		var kapasitas *int
		var distance float64
		if err := rows.Scan(&id, &nama, &kategori, &deskripsi, &kapasitas, &geojson, &distance); err != nil {
			log.Printf("[handler:search] skip baris (scan error): %v", err)
			skipped++
			continue
		}

		// Parse geojson string to json
		var geom map[string]interface{}
		if err := json.Unmarshal([]byte(geojson), &geom); err != nil {
			log.Printf("[handler:search] skip baris (geojson parse error): %v", err)
			skipped++
			continue
		}

		var deskripsiOut any
		if deskripsi != nil {
			deskripsiOut = *deskripsi
		}
		var kapasitasOut any
		if kapasitas != nil {
			kapasitasOut = *kapasitas
		}

		results = append(results, map[string]interface{}{
			"id_potensi":     id,
			"nama_objek":     nama,
			"kategori":       kategori,
			"deskripsi":      deskripsiOut,
			"kapasitas":      kapasitasOut,
			"geometri":       geom,
			"jarak_semantik": distance, // Semakin mendekati 0, maknanya semakin mirip
		})
	}

	// Cek error iterasi yang ditelan diam-diam
	if err := rows.Err(); err != nil {
		return serverError(c, "search.rows", err)
	}
	if skipped > 0 {
		log.Printf("[handler:search] %d baris dilewati (NULL/parse error)", skipped)
	}

	if len(results) == 0 {
		return c.JSON(fiber.Map{"message": "Tidak ada data yang relevan", "data": []interface{}{}})
	}

	return c.JSON(fiber.Map{
		"message": "Pencarian Semantik Berhasil",
		"data":    results,
	})
}
