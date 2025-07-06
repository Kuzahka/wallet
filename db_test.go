package service

import (
	"context"
	_ "database/sql"
	"testing"

	"wallet-api/internal/domain"
	"wallet-api/internal/repository"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Подключение к тестовой БД
	dbURL := "postgres://user:password@localhost:5432/test_wallet_db?sslmode=disable"
	dbPool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer dbPool.Close()

	// Создание репозитория
	repo := repository.NewWalletRepository(dbPool)

	// Создание сервиса
	service := NewWalletService(repo)

	// Тестируем депозит
	walletID, err := domain.ParseUUID("a1b2c3d4-e5f6-7890-g1h2-i3j4k5l6m7n8")
	req := &domain.OperationRequest{
		WalletID:      walletID,
		OperationType: domain.Deposit,
		Amount:        1000,
	}

	err = service.ProcessOperation(context.Background(), req)
	assert.NoError(t, err)

	// Проверяем баланс
	wallet, err := service.GetWalletByID(context.Background(), walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), wallet.Balance)

	// Тестируем списание
	req = &domain.OperationRequest{
		WalletID:      walletID,
		OperationType: domain.Withdraw,
		Amount:        500,
	}

	err = service.ProcessOperation(context.Background(), req)
	assert.NoError(t, err)

	// Проверяем новый баланс
	wallet, err = service.GetWalletByID(context.Background(), walletID)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), wallet.Balance)
}
