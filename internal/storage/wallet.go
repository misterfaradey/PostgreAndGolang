package storage

import (
	"context"
	"github.com/misterfaradey/PostgreAndGolang/internal/dto"
)

type Wallet interface {
	GetWallet(ctx context.Context, id uint64) (dto.Wallet, error)
	GetTransaction(ctx context.Context, id string) (dto.Transaction, error)

	Transfer(ctx context.Context, transaction dto.Transaction) error
}

func (c *db) GetWallet(ctx context.Context, id uint64) (dto.Wallet, error) {

	var wallet dto.Wallet

	err := c.db.QueryRowContext(ctx, getWallet, id).Scan(&wallet.ID, &wallet.Balance)
	return wallet, err
}

func (c *db) GetTransaction(ctx context.Context, id string) (dto.Transaction, error) {

	var t dto.Transaction

	err := c.db.QueryRowContext(ctx, getTransaction, id).Scan(&t.ID, &t.State, &t.Amount)
	return t, err
}

func (c *db) Transfer(ctx context.Context, transaction dto.Transaction) error {

	_, err := c.db.ExecContext(ctx, updateBalance, transaction.ID, transaction.State, transaction.Amount)
	return err
}
