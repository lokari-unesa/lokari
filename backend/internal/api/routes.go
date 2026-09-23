package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lokari/backend/internal/api/handlers"
)

// rateLimiter membatasi request per IP untuk endpoint AI yang memakai kuota berbayar
func rateLimiter(max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak permintaan, coba lagi nanti",
			})
		},
	})
}

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
			"status":  "success",
			"message": "Backend LOKARI is running",
		})
	})

	api.Get("/potensi", h.GetPotensiBencana)

	// AI Routes — dilindungi rate limiter (kuota Cohere/AI berbayar)
	aiLimiter := rateLimiter(30, 1*time.Minute)
	api.Get("/alert", aiLimiter, aiHandler.GetAlert)
	api.Post("/search", aiLimiter, aiHandler.SearchSemantic)

	// News Routes
	api.Get("/news", newsHandler.GetNews)
}
