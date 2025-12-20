package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	PriceService PriceServiceConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type PriceServiceConfig struct {
	APIURL        string
	FetchInterval time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	priceIntervalStr := getEnv("PRICE_FETCH_INTERVAL", "1h")
	priceInterval, err := time.ParseDuration(priceIntervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PRICE_FETCH_INTERVAL: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "assignment"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		PriceService: PriceServiceConfig{
			APIURL:        getEnv("PRICE_API_URL", ""),
			FetchInterval: priceInterval,
		},
	}, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

