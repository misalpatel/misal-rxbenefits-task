// Package util provides utility functions for configuration management.
package util //nolint:revive //Package name is fine IMO

import "os"

// Config holds application configuration. Can be extended to include more
// and work with helm charts.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// InitConfig initializes configuration from environment variables.
func InitConfig() Config {
	return Config{
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     GetEnv("DB_PORT", "5432"),
		DBUser:     GetEnv("DB_USER", "postgres"),
		DBPassword: GetEnv("DB_PASSWORD", "postgres"),
		DBName:     GetEnv("DB_NAME", "dvdrental"),
	}
}

// GetEnv gets an environment variable or returns a default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
