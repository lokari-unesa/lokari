package handlers

import (
	"context"
	"log"
	"time"

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
	if h.DB == nil {
		log.Println("[handler:news] DB pool nil — cek DATABASE_URL di .env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak terhubung"})
	}

	query := `
		SELECT id_kabar, kategori, judul, ringkasan, sumber, created_at 
		FROM kabar_kelud 
		ORDER BY created_at DESC 
		LIMIT 20
	`

	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		return serverError(c, "news.query", err)
	}
	defer rows.Close()

	// Kirim timestamp asli dalam RFC3339 zona WIB, bukan "Baru Saja" hardcoded
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}

	var newsList []NewsResponse
	for rows.Next() {
		var n NewsResponse
		var createdAt time.Time
		if err := rows.Scan(&n.ID, &n.Category, &n.Title, &n.Summary, &n.Source, &createdAt); err != nil {
			log.Printf("[handler:news] skip baris (scan error): %v", err)
			continue
		}
		n.CreatedAt = createdAt.In(loc).Format(time.RFC3339)
		newsList = append(newsList, n)
	}

	// Cek error iterasi yang ditelan diam-diam
	if err := rows.Err(); err != nil {
		return serverError(c, "news.rows", err)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   newsList,
	})
}
