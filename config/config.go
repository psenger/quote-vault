package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBPath     string
	LogLevel   string
	PageSize   int
	CORSOrigin string
}

func Load() *Config {
	pageSize, err := strconv.Atoi(getEnv("PAGE_SIZE", "10"))
	if err != nil {
		pageSize = 10
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "./quotes.db"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		PageSize:   pageSize,
		CORSOrigin: getEnv("CORS_ORIGIN", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}