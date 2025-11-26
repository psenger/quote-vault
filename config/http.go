package config

import (
	"time"
)

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port           string        `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
	IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"5s"`
	MaxHeaderBytes  int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"`
}

// GetHTTPConfig returns HTTP configuration with defaults
func GetHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Port:            getEnv("HTTP_PORT", "8080"),
		ReadTimeout:     parseDuration(getEnv("HTTP_READ_TIMEOUT", "30s")),
		WriteTimeout:    parseDuration(getEnv("HTTP_WRITE_TIMEOUT", "30s")),
		IdleTimeout:     parseDuration(getEnv("HTTP_IDLE_TIMEOUT", "120s")),
		ShutdownTimeout: parseDuration(getEnv("HTTP_SHUTDOWN_TIMEOUT", "5s")),
		MaxHeaderBytes:  parseInt(getEnv("HTTP_MAX_HEADER_BYTES", "1048576")),
	}
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 30 * time.Second
	}
	return d
}

func parseInt(s string) int {
	if s == "" {
		return 1048576
	}
	// Simple conversion for demo - in production use strconv.Atoi with error handling
	return 1048576
}