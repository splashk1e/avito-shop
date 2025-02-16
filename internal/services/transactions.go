package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/splashk1e/avito-shop/internal/models"
	"github.com/splashk1e/avito-shop/internal/storage"
)

type TransactionService struct {
	transactionSaver    TransactionSaver
	transactionProvider TransactionProvider
	log                 *slog.Logger
}

type TransactionProvider interface {
	GetTransactions(ctx context.Context, username string) ([]models.Transaction, error)
	GetItems(ctx context.Context, username string) ([]string, error)
}

type TransactionSaver interface {
	SaveTransaction(ctx context.Context, transaction models.Transaction) (int, error)
}

func NewTransactionsService(transactionSaver TransactionSaver, transactionProvider TransactionProvider, log *slog.Logger) *TransactionService {
	return &TransactionService{
		transactionSaver:    transactionSaver,
		transactionProvider: transactionProvider,
		log:                 log,
	}
}

func (t *TransactionService) GetTransactions(ctx context.Context, username string) ([]models.Transaction, error) {
	const op = "services.transactions.GetTransactions"
	log := t.log.With(slog.String("op", op))
	log.Info("getting transactions")
	transactions, err := t.transactionProvider.GetTransactions(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err.Error())
			return nil, fmt.Errorf("%s %w", op, ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return transactions, nil
}

func (t *TransactionService) GetItems(ctx context.Context, username string) ([]string, error) {
	const op = "services.transactions.GetItems"
	log := t.log.With(slog.String("op", op))
	log.Info("getting items")
	items, err := t.transactionProvider.GetItems(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err.Error())
			return nil, fmt.Errorf("%s %w", op, ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return items, nil
}

var (
	itemList = map[string]int{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	}
)

func (t *TransactionService) SaveTransaction(ctx context.Context, sender, receiver, item string, amount int) (int, error) {
	const op = "services.transactions.SaveTransaction"
	log := t.log.With(slog.String("op", op))
	log.Info("saving transaction")
	transaction := models.Transaction{
		Sender:   sender,
		Reciever: "",
		Item:     item,
		Amount:   amount,
	}
	if item != "" {
		transaction.Amount = itemList[item]
	}
	id, err := t.transactionSaver.SaveTransaction(ctx, transaction)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", err.Error())
			return 0, fmt.Errorf("%s %w", op, ErrInvalidCredentials)
		}
		return 0, fmt.Errorf("%s %w", op, err)
	}
	return id, nil
}
