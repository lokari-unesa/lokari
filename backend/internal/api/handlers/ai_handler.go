package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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
	// 1. Fetch real-time NASA EONET data
	nasaURL := "https://eonet.gsfc.nasa.gov/api/v3/events?category=volcanoes&status=open&limit=5"
	resp, err := http.Get(nasaURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menarik data satelit NASA"})
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI gagal merangkum pesan", "details": err.Error()})
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
	var req SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Input JSON tidak valid"})
	}

	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query tidak boleh kosong"})
	}

	// 1. Convert user's sentence to a vector embedding using HuggingFace
	embedding, err := ai.GenerateEmbedding(req.Query)
	if err != nil {
		if strings.Contains(err.Error(), "KUOTA_COHERE_HABIS") {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "KUOTA_HABIS", "message": "error.quota.search.msg"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat vektor embedding", "details": err.Error()})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengeksekusi pencarian vektor", "details": err.Error()})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, nama, kategori, deskripsi, geojson string
		var kapasitas int
		var distance float64
		if err := rows.Scan(&id, &nama, &kategori, &deskripsi, &kapasitas, &geojson, &distance); err != nil {
			continue
		}
		
		// Parse geojson string to json
		var geom map[string]interface{}
		json.Unmarshal([]byte(geojson), &geom)

		results = append(results, map[string]interface{}{
			"id_potensi": id,
			"nama_objek": nama,
			"kategori":   kategori,
			"deskripsi":  deskripsi,
			"kapasitas":  kapasitas,
			"geometri":   geom,
			"jarak_semantik": distance, // Semakin mendekati 0, maknanya semakin mirip
		})
	}

	if len(results) == 0 {
		return c.JSON(fiber.Map{"message": "Tidak ada data yang relevan", "data": []interface{}{}})
	}

	return c.JSON(fiber.Map{
		"message": "Pencarian Semantik Berhasil",
		"data":    results,
	})
}
