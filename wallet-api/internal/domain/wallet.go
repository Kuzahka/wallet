package domain

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `json:"walletId"`
	Balance int64     `json:"balance"`
}

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type OperationRequest struct {
	WalletID      uuid.UUID     `json:"walletId"`
	OperationType OperationType `json:"operationType"`
	Amount        int64         `json:"amount"`
}

// ParseUUID конвертирует строку в uuid.UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
