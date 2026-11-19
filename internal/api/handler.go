package api

import (
"context"
"log"

"email-delivery-service/internal/db"
"email-delivery-service/internal/models"
"email-delivery-service/internal/queue"

"github.com/gofiber/fiber/v2"
"github.com/google/uuid"
)

type EmailRequest struct {
Recipient string `json:"recipient"`
Subject   string `json:"subject"`
Body      string `json:"body"`
}

func SetupRoutes(app *fiber.App) {
app.Post("/api/emails", func(c *fiber.Ctx) error {
var req EmailRequest
if err := c.BodyParser(&req); err != nil {
return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
}

if req.Recipient == "" || req.Subject == "" {
return c.Status(400).JSON(fiber.Map{"error": "Recipient and Subject are required"})
}

ctx := context.Background()
email := models.Email{
ID:         uuid.New(),
Recipient:  req.Recipient,
Subject:    req.Subject,
Body:       req.Body,
Status:     models.StatusQueued,
RetryCount: 0,
}

query := `INSERT INTO emails (id, recipient, subject, body, status, retry_count) 
          VALUES ($1, $2, $3, $4, $5, $6)`
_, err := db.Pool.Exec(ctx, query, email.ID, email.Recipient, email.Subject, email.Body, email.Status, email.RetryCount)
if err != nil {
log.Printf("DB Error: %v", err)
return c.Status(500).JSON(fiber.Map{"error": "Failed to save email"})
}

if err := queue.EnqueueEmail(ctx, email); err != nil {
log.Printf("Queue Error: %v", err)
return c.Status(500).JSON(fiber.Map{"error": "Failed to queue email"})
}

return c.Status(201).JSON(fiber.Map{
"message": "Email queued successfully",
"id":      email.ID,
})
})

// NEW: Get Delivery History for a specific email
app.Get("/api/emails/:id/history", func(c *fiber.Ctx) error {
emailID := c.Params("id")
ctx := context.Background()

rows, err := db.Pool.Query(ctx, "SELECT provider, status, error_message, attempted_at FROM delivery_attempts WHERE email_id = $1 ORDER BY attempted_at DESC", emailID)
if err != nil {
return c.Status(500).JSON(fiber.Map{"error": "DB Error"})
}
defer rows.Close()

var attempts []models.DeliveryAttempt
for rows.Next() {
var a models.DeliveryAttempt
rows.Scan(&a.Provider, &a.Status, &a.ErrorMsg, &a.AttemptedAt)
attempts = append(attempts, a)
}

return c.JSON(attempts)
})

// NEW: Delivery Analytics
app.Get("/api/analytics", func(c *fiber.Ctx) error {
ctx := context.Background()
var delivered, failed, queued int

db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM emails WHERE status = 'DELIVERED'").Scan(&delivered)
db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM emails WHERE status = 'FAILED'").Scan(&failed)
db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM emails WHERE status = 'QUEUED' OR status = 'PROCESSING'").Scan(&queued)

return c.JSON(fiber.Map{
"delivered": delivered,
"failed":    failed,
"in_queue":  queued,
"total":     delivered + failed + queued,
})
})
}