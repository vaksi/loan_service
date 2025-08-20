package config

import (
    "fmt"
    "os"
)

// Config holds configuration values for the application.
// Each field can be configured using environment variables prefixed with
// the field name (for example DB_HOST, DB_PORT, etc.).
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    ServerPort string
}

// Load reads configuration from environment variables and sets default
// values when variables are not present. These defaults work well for
// local development using docker-compose where the Postgres database
// service is named `db` and exposes port 5432.
func Load() Config {
    cfg := Config{
        DBHost:     getEnv("DB_HOST", "db"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "amartha"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
    }
    return cfg
}

// DSN returns the Postgres Data Source Name constructed from the
// configuration. This DSN is used by GORM to open a database
// connection.
func (c Config) DSN() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// getEnv returns the value of the given environment variable if it
// exists, otherwise it returns the provided default value. This helper
// function centralizes environment variable lookups and avoids
// repetitive code throughout the application.
func getEnv(key, defaultVal string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultVal
}