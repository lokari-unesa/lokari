package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/lokari/backend/internal/api"
	"github.com/lokari/backend/internal/config"
	"github.com/lokari/backend/internal/worker"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to Database
	var db *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
		} else {
			db = pool
			defer db.Close()
			log.Println("Connected to PostgreSQL successfully.")
		}
	} else {
		log.Println("DATABASE_URL is not set, skipping database connection.")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "LOKARI API",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Pass dependencies to API routes
	api.SetupRoutes(app, db)

	// Start Zero-Admin Cron Jobs
	cronJob := worker.StartCronJobs(db)
	defer cronJob.Stop()

	// Start server
	log.Printf("Starting server dynamically on port %s (loaded from .env)", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
