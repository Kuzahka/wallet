package api_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-api/internal/api"
	"wallet-api/internal/domain"
	_ "wallet-api/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для WalletService
type mockWalletService struct {
	mock.Mock
}

func (m *mockWalletService) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *mockWalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	args := m.Called(ctx, id)
	if wallet, ok := args.Get(0).(*domain.Wallet); ok {
		return wallet, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockWalletService) ProcessOperation(ctx context.Context, req *domain.OperationRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestAPI_Deposit(t *testing.T) {
	mockService := &mockWalletService{}
	mockService.On("ProcessOperation", mock.Anything, mock.Anything).Return(nil)

	apiHandler := api.NewWalletAPI(mockService)
	server := httptest.NewServer(apiHandler.Routes())
	defer server.Close()

	payload := []byte(`{
        "walletId": "a81bc81b-dead-4e5d-abff-90865d1e13b1",
        "operationType": "DEPOSIT",
        "amount": 1000
    }`)

	resp, err := http.Post(server.URL+"/api/v1/wallet", "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAPI_Withdraw_NotEnoughBalance(t *testing.T) {
	mockService := &mockWalletService{}
	mockService.On("ProcessOperation", mock.Anything, mock.Anything).Return(errors.New("not enough balance"))

	apiHandler := api.NewWalletAPI(mockService)
	server := httptest.NewServer(apiHandler.Routes())
	defer server.Close()

	payload := []byte(`{
        "walletId": "a81bc81b-dead-4e5d-abff-90865d1e13b1",
        "operationType": "WITHDRAW",
        "amount": 1000
    }`)

	resp, err := http.Post(server.URL+"/api/v1/wallet", "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
