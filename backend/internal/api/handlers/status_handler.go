package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// stateKeyKeludStatus sama dengan kunci yang dipakai service.FetchKeludStatus
// (backend/internal/service/kelud.go) — nilai berbentuk "Level I (Normal)" dll.
const stateKeyKeludStatus = "kelud_status"

type StatusHandler struct {
	DB *pgxpool.Pool
}

func NewStatusHandler(db *pgxpool.Pool) *StatusHandler {
	return &StatusHandler{DB: db}
}

// GetKeludStatus mengembalikan status tingkat aktivitas Gunung Kelud terakhir
// yang disimpan monitor_state oleh job pemantauan (hot loop). Frontend memakai
// endpoint ini sebagai satu-satunya sumber kebenaran untuk kartu status di
// halaman utama — bukan menebak dari berita teratas.
func (h *StatusHandler) GetKeludStatus(c *fiber.Ctx) error {
	if h.DB == nil {
		log.Println("[handler:status] DB pool nil — cek DATABASE_URL di .env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak terhubung"})
	}

	var value string
	var updatedAt time.Time
	err := h.DB.QueryRow(context.Background(),
		`SELECT value, updated_at FROM monitor_state WHERE key = $1`, stateKeyKeludStatus).
		Scan(&value, &updatedAt)
	if err == pgx.ErrNoRows {
		// Belum pernah tercatat status — biarkan frontend menampilkan state
		// netral, bukan error.
		return c.JSON(fiber.Map{"status": "success", "data": nil})
	}
	if err != nil {
		return serverError(c, "status.query", err)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"status":     value,
			"updated_at": updatedAt.In(loc).Format(time.RFC3339),
		},
	})
}