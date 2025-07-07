package domain

import (
	"errors"
)

var (
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrInvalidUUID      = errors.New("invalid wallet ID")
	ErrUnknownOp        = errors.New("unknown operation type")
)
