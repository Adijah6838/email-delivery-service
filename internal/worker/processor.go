package worker

import (
"context"
"encoding/json"
"log"
"math/rand"
"time"

"email-delivery-service/internal/db"
"email-delivery-service/internal/models"
"email-delivery-service/internal/queue"
)

const MaxRetries = 3

func StartWorker(ctx context.Context, workerID int) {
log.Printf("Worker %d started...", workerID)
for {
// BRPOP waits for an item (timeout 0 = wait forever)
result, err := db.RedisClient.BRPop(ctx, 0, queue.EmailQueueKey).Result()
if err != nil {
log.Printf("Worker %d queue error: %v", workerID, err)
continue
}

var email models.Email
if err := json.Unmarshal([]byte(result[1]), &email); err != nil {
log.Printf("Worker %d unmarshal error: %v", workerID, err)
continue
}

log.Printf("Worker %d processing email %s (Attempt %d)", workerID, email.ID, email.RetryCount+1)

// 1. Update status to PROCESSING
db.Pool.Exec(ctx, "UPDATE emails SET status = $1 WHERE id = $2", models.StatusProcessing, email.ID)

// 2. Simulate Provider Delay
time.Sleep(500 * time.Millisecond)

// 3. Simulate Success/Failure (50% failure rate for testing)
if rand.Float32() < 0.5 {
// SUCCESS
db.Pool.Exec(ctx, "UPDATE emails SET status = $1 WHERE id = $2", models.StatusDelivered, email.ID)
db.Pool.Exec(ctx, "INSERT INTO delivery_attempts (email_id, provider, status) VALUES ($1, $2, $3)", email.ID, "MockProvider", "SUCCESS")
log.Printf("Worker %d SUCCESS: Email %s delivered", workerID, email.ID)
} else {
// FAILURE
email.RetryCount++

if email.RetryCount < MaxRetries {
// RETRY: Push back to Redis and update DB
db.Pool.Exec(ctx, "UPDATE emails SET status = $1, retry_count = $2 WHERE id = $3", models.StatusQueued, email.RetryCount, email.ID)
db.Pool.Exec(ctx, "INSERT INTO delivery_attempts (email_id, provider, status, error_message) VALUES ($1, $2, $3, $4)", email.ID, "MockProvider", "FAILED", "Transient error, retrying")

// Re-enqueue immediately for demo (in production, use a delayed queue)
queue.EnqueueEmail(ctx, email)
log.Printf("Worker %d RETRY: Email %s failed, re-queued (Attempt %d)", workerID, email.ID, email.RetryCount)
} else {
// DEAD LETTER QUEUE (Permanent Failure)
db.Pool.Exec(ctx, "UPDATE emails SET status = $1, retry_count = $2 WHERE id = $3", models.StatusFailed, email.RetryCount, email.ID)
db.Pool.Exec(ctx, "INSERT INTO delivery_attempts (email_id, provider, status, error_message) VALUES ($1, $2, $3, $4)", email.ID, "MockProvider", "FAILED", "Max retries exceeded")
log.Printf("Worker %d DLQ: Email %s permanently failed after %d attempts", workerID, email.ID, email.RetryCount)
}
}
}
}