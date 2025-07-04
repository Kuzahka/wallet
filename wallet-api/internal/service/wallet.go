package service

import (
	"context"
	"errors"
	"wallet-api/internal/domain"
	"wallet-api/internal/repository"

	"github.com/google/uuid"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	return s.repo.Create(ctx, wallet)
}

func (s *WalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WalletService) ProcessOperation(ctx context.Context, req *domain.OperationRequest) error {
	_, err := s.repo.GetByID(ctx, req.WalletID)
	if err != nil {
		return err
	}

	switch req.OperationType {
	case domain.Deposit:
		err = s.repo.UpdateBalance(ctx, req.WalletID, req.Amount)
	case domain.Withdraw:
		err = s.repo.UpdateBalance(ctx, req.WalletID, -req.Amount)
	default:
		return errors.New("invalid operation type")
	}

	return err
}
