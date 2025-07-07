package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"wallet-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error)
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error
}

type walletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) WalletRepository {
	return &walletRepository{db: db}
}

// Create создаёт новый кошелёк в БД
func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `
        INSERT INTO wallets (id, balance)
        VALUES ($1, $2)
        ON CONFLICT (id) DO NOTHING`

	_, err := r.db.Exec(ctx, query, wallet.ID, wallet.Balance)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

// GetByID возвращает кошелёк по ID
func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	row := r.db.QueryRow(ctx, "SELECT id, balance FROM wallets WHERE id = $1", id)

	var wallet domain.Wallet
	err := row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}

	return &wallet, nil
}

// UpdateBalance изменяет баланс кошелька атомарно с блокировкой строки
func (r *walletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Блокируем строку для последовательной обработки
	var currentBalance sql.NullInt64
	err = tx.QueryRow(ctx, "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", walletID).Scan(&currentBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrWalletNotFound
		}
		return fmt.Errorf("failed to lock wallet row: %w", err)
	}

	if !currentBalance.Valid {
		return fmt.Errorf("balance is NULL for wallet ID: %s", walletID)
	}

	// Обновляем баланс
	_, err = tx.Exec(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Коммитим транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
