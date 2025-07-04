package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// InitSchema создаёт таблицы, если их нет в БД
func InitSchema(ctx context.Context, db *pgxpool.Pool) error {
	query := `
    CREATE TABLE IF NOT EXISTS wallets (
        id UUID PRIMARY KEY,
        balance BIGINT NOT NULL DEFAULT 0
    );`

	_, err := db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}
