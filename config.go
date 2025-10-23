// Файл: config/config.go
package config

import (
	"fmt"
	"os"
)

// Config хранит конфигурацию приложения.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load загружает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "host.docker.internal"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := "warehouse_user"
	if dbUser == "" {
		return nil, fmt.Errorf("переменная окружения DB_USER не установлена")
	}

	dbPassword := "postgres"
	if dbPassword == "" {
		return nil, fmt.Errorf("переменная окружения DB_PASSWORD не установлена")
	}

	dbName := "warehouse_db"
	if dbName == "" {
		return nil, fmt.Errorf("переменная окружения DB_NAME не установлена")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	return &Config{
		DatabaseURL: dsn,
		Port:        port,
	}, nil
}
