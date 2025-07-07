package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	App struct {
		Port string
	}
}

func LoadConfig() Config {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println("No config.env file found or error loading it")
	}

	return Config{
		DB: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
			SSLMode  string
		}{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "wallet_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		App: struct {
			Port string
		}{
			Port: getEnv("PORT", "8080"),
		},
	}
}

// getEnv — вспомогательная функция для получения переменной окружения или значения по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
