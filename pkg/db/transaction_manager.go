package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TransactionManager interface {
	BeginTransaction(ctx context.Context) (*sqlx.Tx, error)
}

type transactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) BeginTransaction(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := tm.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
