package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all environment configuration for the application.
// Node analogy: this is like your config.js that does dotenv.config()
// and exports a typed config object.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
	AWS      AWSConfig
}

type AppConfig struct {
	Env  string // "development" | "production" | "test"
	Port string // "8080"
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN builds the PostgreSQL connection string.
// Node analogy: like building your DATABASE_URL string for pg or prisma.
// DSN = Data Source Name — standard term across all languages.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type AWSConfig struct {
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
}

// Load reads the .env file and returns a populated Config struct.
// Node analogy: this is your one-time dotenv.config() + exporting config object.
// Call this ONCE in main.go and pass it down everywhere.
func Load(envFile string) (*Config, error) {
	// Load .env file — in production (Docker), env vars are injected directly
	// so this won't fail. The _ means "ignore error if file not found."
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("No env file found at %s, using system env vars\n", envFile)
		}
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "job_automator"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", ""),
			ExpiryHours: jwtExpiry,
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			GroupID: getEnv("KAFKA_GROUP_ID", "job-automator-group"),
		},
		AWS: AWSConfig{
			Region:     getEnv("AWS_REGION", "ap-south-1"),
			BucketName: getEnv("AWS_BUCKET_NAME", ""),
			AccessKey:  getEnv("AWS_ACCESS_KEY", ""),
			SecretKey:  getEnv("AWS_SECRET_KEY", ""),
		},
	}

	// Fail fast — if critical config is missing, crash immediately.
	// Better to crash at startup than fail mysteriously at runtime.
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required but not set")
	}

	return cfg, nil
}

// getEnv is a helper: read env var, fall back to default if not set.
// Node analogy: process.env.PORT || '8080'
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
