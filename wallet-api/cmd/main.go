package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"wallet-api/config"
	"wallet-api/internal/api"
	"wallet-api/internal/repository"
	"wallet-api/internal/service"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	cfg := config.LoadConfig()

	// Формируем DSN из конфига
	dbDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	// Настройка пула соединений
	poolConfig, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}
	poolConfig.MaxConns = 10

	// Подключение к PostgreSQL
	dbPool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Инициализация схемы
	if err := repository.InitSchema(context.Background(), dbPool); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Сервисы и API
	repo := repository.NewWalletRepository(dbPool)
	service := service.NewWalletService(repo)
	apiHandler := api.NewWalletAPI(service)

	// Запуск сервера
	log.Printf("Server is running on port %s\n", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, apiHandler.Routes()))
}
