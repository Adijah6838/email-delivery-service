package config

import (
"log"
"os"
"path/filepath"

"github.com/joho/godotenv"
)

type Config struct {
Port        string
DatabaseURL string
RedisURL    string
}

func LoadConfig() *Config {
// Try to load .env from multiple possible locations
envPaths := []string{
".env",
filepath.Join("..", ".env"),
filepath.Join("..", "..", ".env"),
}

loaded := false
for _, path := range envPaths {
if err := godotenv.Load(path); err == nil {
log.Printf("Loaded .env from %s", path)
loaded = true
break
}
}

if !loaded {
log.Println("No .env file found, using system environment variables")
}

return &Config{
Port:        getEnv("PORT", "8080"),
DatabaseURL: getEnv("DATABASE_URL", ""),
RedisURL:    getEnv("REDIS_URL", ""),
}
}

func getEnv(key, fallback string) string {
if value, exists := os.LookupEnv(key); exists {
return value
}
return fallback
}