package service

import (
	"context"
	_ "errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"wallet-api/internal/domain"
	_ "wallet-api/internal/repository"
)

// Мок репозитория
type mockWalletRepository struct {
	mock.Mock
}

func (m *mockWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *mockWalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	args := m.Called(ctx, id)
	if wallet, ok := args.Get(0).(*domain.Wallet); ok {
		return wallet, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockWalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	args := m.Called(ctx, walletID, amount)
	return args.Error(0)
}

func TestProcessOperation_Deposit(t *testing.T) {
	mockRepo := new(mockWalletRepository)
	service := NewWalletService(mockRepo)

	walletID, _ := uuid.Parse("a1b2c3d4-e5f6-4789-90ab-cdef12345678")

	// Мок GetByID
	mockRepo.On("GetByID", mock.Anything, walletID).Return(&domain.Wallet{
		ID:      walletID,
		Balance: 1000,
	}, nil)

	// Мок UpdateBalance
	mockRepo.On("UpdateBalance", mock.Anything, walletID, int64(500)).Return(nil)

	req := &domain.OperationRequest{
		WalletID:      walletID,
		OperationType: domain.Deposit,
		Amount:        500,
	}

	err := service.ProcessOperation(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessOperation_Withdraw_Success(t *testing.T) {
	mockRepo := new(mockWalletRepository)
	service := NewWalletService(mockRepo)

	walletID, _ := uuid.Parse("a1b2c3d4-e5f6-4789-90ab-cdef12345678")

	mockRepo.On("GetByID", mock.Anything, walletID).Return(&domain.Wallet{
		ID:      walletID,
		Balance: 1000,
	}, nil)

	mockRepo.On("UpdateBalance", mock.Anything, walletID, int64(-500)).Return(nil)

	req := &domain.OperationRequest{
		WalletID:      walletID,
		OperationType: domain.Withdraw,
		Amount:        500,
	}

	err := service.ProcessOperation(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessOperation_Withdraw_NotEnoughBalance(t *testing.T) {
	mockRepo := new(mockWalletRepository)
	service := NewWalletService(mockRepo)

	walletID, _ := uuid.Parse("a1b2c3d4-e5f6-4789-90ab-cdef12345678")

	// Мок GetByID возвращает кошелёк с балансом 100
	mockRepo.On("GetByID", mock.Anything, walletID).Return(&domain.Wallet{
		ID:      walletID,
		Balance: 100,
	}, nil)

	// Убедись, что UpdateBalance НЕ вызывается
	mockRepo.On("UpdateBalance", mock.Anything, walletID, int64(-200)).Maybe().Return(nil)

	req := &domain.OperationRequest{
		WalletID:      walletID,
		OperationType: domain.Withdraw,
		Amount:        200,
	}

	err := service.ProcessOperation(context.Background(), req)
	assert.Error(t, err)
	assert.EqualError(t, err, "not enough balance")
	mockRepo.AssertExpectations(t)
}

func TestProcessOperation_InvalidUUID(t *testing.T) {
	mockRepo := new(mockWalletRepository)
	service := NewWalletService(mockRepo)

	req := &domain.OperationRequest{
		WalletID:      uuid.Nil,
		OperationType: domain.Deposit,
		Amount:        500,
	}

	err := service.ProcessOperation(context.Background(), req)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid wallet ID")
	mockRepo.AssertExpectations(t)
}
