package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/splashk1e/avito-shop/internal/models"
	"github.com/splashk1e/avito-shop/internal/storage"
)

type Strorage struct {
	pool *pgxpool.Pool
	mu   sync.RWMutex
}

const (
	saveUserQuery        = "INSERT INTO users(username, pass_hash) VALUES($1, $2) RETURNING id"
	getUserQuery         = "SELECT id, username, pass_hash FROM users WHERE username = $1"
	saveTransactionQuery = "INSERT INTO transactions(sender, reciever, amount, item) VALUES($1, $2, $3, $4) RETURNING id"
	getTransactionsQuery = "SELECT id, sender, reciever, amount, item FROM transactions WHERE sender = $1 OR reciever = $1"
	getItemsQuery        = "SELECT item FROM transactions WHERE sender = $1 AND reciever=NULL"
)

func New(ctx context.Context, storagePath string) *Strorage {
	pool, err := pgxpool.New(ctx, storagePath)
	if err != nil {
		panic("can't connect to postgres with error:" + err.Error())
	}
	if err := pool.Ping(ctx); err != nil {
		panic("postgres don't answer with error:" + err.Error())
	}
	return &Strorage{
		pool: pool,
		mu:   sync.RWMutex{},
	}
}

func (s *Strorage) SaveUser(ctx context.Context, username string, passHash []byte) (int, error) {
	const op = "postgres.storage.SaveUser"
	s.mu.Lock()
	defer s.mu.Unlock()
	var id int
	if user, _ := s.GetUser(ctx, username); user != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}
	row, err := s.pool.Query(ctx, saveUserQuery, username, passHash)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer row.Conn()
	if !row.Next() {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Strorage) GetUser(ctx context.Context, username string) (*models.User, error) {
	const op = "postgres.storage.GetUser"
	var user models.User
	s.mu.RLock()
	defer s.mu.RUnlock()
	row, err := s.pool.Query(ctx, getUserQuery, username)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !row.Next() {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}
	if err := row.Scan(&user.Id, &user.Username, &user.PassHash); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}
