package handlers

import (
	"context"
	"database/sql"
	"log"

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

// serverError mencatat error asli di log server dan mengembalikan pesan generik
// ke client agar detail internal/skema DB tidak bocor.
func serverError(c *fiber.Ctx, op string, err error) error {
	log.Printf("[handler:%s] error: %v", op, err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Terjadi kesalahan internal pada server",
	})
}

// GetPotensiBencana fetches all spatial data from potensi_bencana table
func (h *Handler) GetPotensiBencana(c *fiber.Ctx) error {
	if h.DB == nil {
		log.Println("[handler:potensi] DB pool nil — cek DATABASE_URL di .env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database tidak terhubung"})
	}

	// Fetch with ST_AsGeoJSON to get geometry as JSON string
	query := `
		SELECT id_potensi, nama_objek, kategori, tingkat_risiko, deskripsi, alamat_dusun, 
		       kapasitas_orang, ST_AsGeoJSON(geometri) as geometri, foto_lokasi, kontak_darurat, created_at, updated_at
		FROM potensi_bencana
	`
	rows, err := h.DB.Query(context.Background(), query)
	if err != nil {
		return serverError(c, "potensi.query", err)
	}
	defer rows.Close()

	var results []models.PotensiBencana
	var skipped int
	for rows.Next() {
		var p models.PotensiBencana
		var geometri sql.NullString
		err := rows.Scan(
			&p.IDPotensi, &p.NamaObjek, &p.Kategori, &p.TingkatRisiko, &p.Deskripsi,
			&p.AlamatDusun, &p.KapasitasOrang, &geometri, &p.FotoLokasi, &p.KontakDarurat,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			log.Printf("[handler:potensi] skip baris (scan error): %v", err)
			skipped++
			continue
		}
		// GeoJSON NULL (baris tanpa geometri) → lewati baris, jangan gagalkan seluruh request
		if !geometri.Valid || geometri.String == "" {
			log.Printf("[handler:potensi] skip baris tanpa geometri: %s", p.IDPotensi)
			skipped++
			continue
		}
		p.Geometri = geometri.String
		results = append(results, p)
	}

	// Cek error iterasi yang ditelan diam-diam
	if err := rows.Err(); err != nil {
		return serverError(c, "potensi.rows", err)
	}
	if skipped > 0 {
		log.Printf("[handler:potensi] %d baris dilewati (NULL/scan error)", skipped)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   results,
	})
}
