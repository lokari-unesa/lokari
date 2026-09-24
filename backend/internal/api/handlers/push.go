package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PushHandler struct {
	db *pgxpool.Pool
}

func NewPushHandler(db *pgxpool.Pool) *PushHandler {
	return &PushHandler{db: db}
}

// Subscribe saves a Web Push Subscription to the database
func (h *PushHandler) Subscribe(c *fiber.Ctx) error {
	type Keys struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	}
	type Payload struct {
		Endpoint string `json:"endpoint"`
		Keys     Keys   `json:"keys"`
	}

	var payload Payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if payload.Endpoint == "" || payload.Keys.P256dh == "" || payload.Keys.Auth == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing push subscription data"})
	}

	_, err := h.db.Exec(context.Background(), 
		"INSERT INTO push_subscriptions (endpoint, p256dh, auth) VALUES ($1, $2, $3) ON CONFLICT (endpoint) DO NOTHING",
		payload.Endpoint, payload.Keys.P256dh, payload.Keys.Auth,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save subscription"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Subscription saved"})
}
