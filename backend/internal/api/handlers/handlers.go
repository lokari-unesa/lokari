package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/models"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

// GetPotensiBencana fetches all spatial data from potensi_bencana table
func (h *Handler) GetPotensiBencana(c *fiber.Ctx) error {
	if h.DB == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database not connected"})
	}

	// Fetch with ST_AsGeoJSON to get geometry as JSON string
	query := `
		SELECT id_potensi, nama_objek, kategori, tingkat_risiko, deskripsi, alamat_dusun, 
		       kapasitas_orang, ST_AsGeoJSON(geometri) as geometri, foto_lokasi, kontak_darurat, created_at, updated_at
		FROM potensi_bencana
	`
	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var results []models.PotensiBencana
	for rows.Next() {
		var p models.PotensiBencana
		err := rows.Scan(
			&p.IDPotensi, &p.NamaObjek, &p.Kategori, &p.TingkatRisiko, &p.Deskripsi,
			&p.AlamatDusun, &p.KapasitasOrang, &p.Geometri, &p.FotoLokasi, &p.KontakDarurat,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		results = append(results, p)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   results,
	})
}
