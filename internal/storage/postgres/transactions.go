package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/splashk1e/avito-shop/internal/models"
	"github.com/splashk1e/avito-shop/internal/storage"
)

func (s *Strorage) GetTransactions(ctx context.Context, username string) ([]models.Transaction, error) {
	const op = "postgres.storage.GetTransactions"
	transactions := make([]models.Transaction, 0)
	var transaction models.Transaction
	s.mu.RLock()
	defer s.mu.RUnlock()
	rows, err := s.pool.Query(ctx, getTransactionsQuery, username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		if err := rows.Scan(&transaction.Id, &transaction.Sender, &transaction.Reciever, &transaction.Amount, &transaction.Item); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, transaction)
	}
	if len(transactions) == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}
	return transactions, nil
}

func (s *Strorage) GetItems(ctx context.Context, username string) ([]string, error) {
	const op = "postgres.storage.GetItems"
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]string, 0)
	var item string
	rows, err := s.pool.Query(ctx, getItemsQuery, username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Strorage) SaveTransaction(ctx context.Context, transaction models.Transaction) (int, error) {
	const op = "postgres.storage.SaveTransactions"
	s.mu.Lock()
	defer s.mu.Unlock()

	var id int
	user, err := s.GetUser(ctx, transaction.Reciever)
	if errors.Is(err, storage.ErrUserNotFound) && transaction.Item != "" {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}
	if transaction.Amount > user.Coins {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrNoCoins)
	}
	row, err := s.pool.Query(ctx, saveTransactionQuery, transaction.Sender, transaction.Reciever, transaction.Amount, transaction.Item)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if row.Next() {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
