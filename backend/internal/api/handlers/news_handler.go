package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsHandler struct {
	DB *pgxpool.Pool
}

func NewNewsHandler(db *pgxpool.Pool) *NewsHandler {
	return &NewsHandler{DB: db}
}

type NewsResponse struct {
	ID        string `json:"id"`
	Category  string `json:"category"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
}

func (h *NewsHandler) GetNews(c *fiber.Ctx) error {
	query := `
		SELECT id_kabar, kategori, judul, ringkasan, sumber, created_at 
		FROM kabar_kelud 
		ORDER BY created_at DESC 
		LIMIT 20
	`

	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Gagal menarik berita: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menarik berita dari database"})
	}
	defer rows.Close()

	var newsList []NewsResponse
	for rows.Next() {
		var n NewsResponse
		var createdAt interface{}
		err := rows.Scan(&n.ID, &n.Category, &n.Title, &n.Summary, &n.Source, &createdAt)
		if err != nil {
			log.Printf("Error scanning news row: %v\n", err)
			continue
		}
		
		// Simplify date handling for now
		n.CreatedAt = "Baru Saja"
		
		newsList = append(newsList, n)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   newsList,
	})
}
