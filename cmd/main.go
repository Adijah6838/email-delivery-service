package main

import (
"context"
"log"

"email-delivery-service/internal/api"
"email-delivery-service/internal/db"
"email-delivery-service/internal/worker"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
ctx := context.Background()

dbURL := "postgres://email_user:email_password@localhost:5450/email_db"
redisURL := "redis://localhost:6390/0"

if err := db.ConnectDB(ctx, dbURL); err != nil {
log.Fatalf("Failed to connect to PostgreSQL: %v", err)
}

if err := db.ConnectRedis(ctx, redisURL); err != nil {
log.Fatalf("Failed to connect to Redis: %v", err)
}

// 1. Start 3 Background Workers
for i := 1; i <= 3; i++ {
go worker.StartWorker(ctx, i)
}

// 2. Start API Server
app := fiber.New()
app.Use(logger.New())
api.SetupRoutes(app)

log.Println("Starting API server on http://localhost:8090")
if err := app.Listen(":8090"); err != nil {
log.Fatalf("Server error: %v", err)
}
}