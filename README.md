# Advanced Email Delivery Simulation Service

A high-concurrency backend service built with Go, Redis, and PostgreSQL that simulates an email delivery system with background workers, retry policies, and a Dead Letter Queue.

## Features
- **REST API**: Accept email requests and queue them for processing.
- **Message Queue**: Uses Redis `BRPOP` for reliable, concurrent job processing.
- **Background Workers**: 3 concurrent Go routines process emails asynchronously.
- **Provider Simulation**: Simulates a 50% failure rate from a mock email provider.
- **Retry Policy**: Failed emails are re-queued up to 3 times (Max Retries).
- **Dead Letter Queue (DLQ)**: Emails that fail 3 times are permanently marked as `FAILED`.
- **Delivery History**: Tracks every single attempt (success or failure) for each email.
- **Analytics**: Real-time stats on delivered, failed, and queued emails.
- **Structured Logging**: Every worker action is logged with the worker ID and email ID.

## Tech Stack
- **Language**: Go (Golang) 1.21+
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL 15
- **Queue**: Redis 7
- **Infrastructure**: Docker & Docker Compose

## Quick Start
1. `docker-compose up -d`
2. `Get-Content migrations/001_init.sql | docker exec -i email_postgres psql -U email_user -d email_db`
3. `go run cmd/main.go`

## API Endpoints
- `POST /api/emails`: Submit a new email for delivery.
- `GET /api/analytics`: Get system-wide delivery stats.
- `GET /api/emails/{id}/history`: Get the delivery attempts for a specific email.