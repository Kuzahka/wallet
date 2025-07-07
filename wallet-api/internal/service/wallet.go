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

type WalletServiceInterface interface {
	CreateWallet(ctx context.Context, wallet *domain.Wallet) error
	GetWalletByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	ProcessOperation(ctx context.Context, req *domain.OperationRequest) error
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
	if req.WalletID == uuid.Nil {
		return errors.New("invalid wallet ID")
	}

	wallet, err := s.repo.GetByID(ctx, req.WalletID)
	if err != nil {
		return err
	}

	switch req.OperationType {
	case domain.Deposit:
		return s.repo.UpdateBalance(ctx, wallet.ID, req.Amount)
	case domain.Withdraw:
		if wallet.Balance < req.Amount {
			return domain.ErrNotEnoughBalance
		}
		return s.repo.UpdateBalance(ctx, wallet.ID, -req.Amount)
	default:
		return domain.ErrUnknownOp
	}
}
