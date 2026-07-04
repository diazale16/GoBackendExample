package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	S3Endpoint  string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3ProjectID string
	ServerPort  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	stage := getEnv("STAGE", "testing")

	var prefix string
	switch stage {
	case "production":
		prefix = "SUPABASE_PRODUCTION_"
	default:
		prefix = "SUPABASE_TESTING_"
	}

	cfg := &Config{
		DatabaseURL: getEnv(prefix+"DATABASE_URL", ""),
		S3Endpoint:  getEnv(prefix+"S3_ENDPOINT", ""),
		S3Region:    getEnv(prefix+"S3_REGION", "us-east-1"),
		S3AccessKey: getEnv(prefix+"S3_ACCESS_ID", ""),
		S3SecretKey: getEnv(prefix+"S3_SECRET_KEY", ""),
		S3Bucket:    getEnv(prefix+"S3_BUCKET", "test-bucket"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	cfg.S3ProjectID = extractProjectID(cfg.S3Endpoint)

	log.Printf("Stage: %s, ProjectID: %s", stage, cfg.S3ProjectID)

	return cfg
}

func extractProjectID(endpoint string) string {
	// endpoint format: https://<project>.storage.supabase.co/storage/v1/s3
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	parts := strings.SplitN(endpoint, ".", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return c.DatabaseURL
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%s", c.ServerPort)
}