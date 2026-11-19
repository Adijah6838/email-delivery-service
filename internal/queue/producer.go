package queue

import (
"context"
"encoding/json"
"log"

"email-delivery-service/internal/db"
"email-delivery-service/internal/models"
)

const EmailQueueKey = "email:queue"

func EnqueueEmail(ctx context.Context, email models.Email) error {
// Convert the Email struct to JSON
data, err := json.Marshal(email)
if err != nil {
return err
}

// Push to the left of the Redis list
if err := db.RedisClient.LPush(ctx, EmailQueueKey, data).Err(); err != nil {
return err
}

log.Printf("Enqueued email %s for %s", email.ID, email.Recipient)
return nil
}