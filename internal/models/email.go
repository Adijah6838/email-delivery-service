package models

import (
"time"
"github.com/google/uuid"
)

type EmailStatus string

const (
StatusQueued     EmailStatus = "QUEUED"
StatusProcessing EmailStatus = "PROCESSING"
StatusDelivered  EmailStatus = "DELIVERED"
StatusFailed     EmailStatus = "FAILED"
)

type Email struct {
ID          uuid.UUID   `json:"id"`
Recipient   string      `json:"recipient"`
Subject     string      `json:"subject"`
Body        string      `json:"body"`
Status      EmailStatus `json:"status"`
RetryCount  int         `json:"retry_count"`
CreatedAt   time.Time   `json:"created_at"`
UpdatedAt   time.Time   `json:"updated_at"`
}

type DeliveryAttempt struct {
ID         uuid.UUID `json:"id"`
EmailID    uuid.UUID `json:"email_id"`
Provider   string    `json:"provider"`
Status     string    `json:"status"`
ErrorMsg   string    `json:"error_message,omitempty"`
AttemptedAt time.Time `json:"attempted_at"`
}