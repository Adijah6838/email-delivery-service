package db

import (
"context"
"log"

"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB(ctx context.Context, dbURL string) error {
var err error
Pool, err = pgxpool.New(ctx, dbURL)
if err != nil {
return err
}

if err := Pool.Ping(ctx); err != nil {
return err
}

log.Println("Connected to PostgreSQL")
return nil
}