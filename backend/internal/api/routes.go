package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/api/handlers"
)

func SetupRoutes(app *fiber.App, db *pgxpool.Pool) {
	// Initialize handlers
	h := handlers.NewHandler(db)
	aiHandler := handlers.NewAIHandler(db)
	newsHandler := handlers.NewNewsHandler(db)

	// Root route to show a welcome message instead of 404
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("LOKARI Backend API is running! Access /api/health to check status.")
	})

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "success",
			"message": "Backend LOKARI is running",
		})
	})

	api.Get("/potensi", h.GetPotensiBencana)

	// AI Routes
	api.Get("/alert", aiHandler.GetAlert)
	api.Post("/search", aiHandler.SearchSemantic)

	// News Routes
	api.Get("/news", newsHandler.GetNews)
}
