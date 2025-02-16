package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/splashk1e/avito-shop/internal/lib/jwt"
	"github.com/splashk1e/avito-shop/internal/models"
	"github.com/splashk1e/avito-shop/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userSaver        UserSaver
	userProvider     UserProvider
	transactionSaver TransactionSaver
	log              *slog.Logger
	TokenTTL         time.Duration
	secret           string
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, passHash []byte) (int, error)
}

type UserProvider interface {
	GetUser(ctx context.Context, username string) (*models.User, error)
}

func NewAuthService(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, TokenTTL time.Duration, secret string) *AuthService {
	return &AuthService{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		TokenTTL:     TokenTTL,
		secret:       secret,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	const op = "services.auth.Login"
	log := a.log.With(slog.String("op", op))
	log.Info("Logining user")
	user, err := a.userProvider.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			if _, err := a.Register(ctx, username, password); err != nil {
				return "", fmt.Errorf("%s %w", op, err)
			}
		}
		return "", fmt.Errorf("%s %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		log.Info("invalid password", err.Error())
		return "", fmt.Errorf("%s %w", op, ErrInvalidCredentials)
	}
	token, err := jwt.NewToken(*user, a.TokenTTL, a.secret)
	if err != nil {
		log.Error("failed to generate token", err.Error())
		return "", fmt.Errorf("%s %w", op, err)
	}
	return token, nil
}

func (a *AuthService) Register(ctx context.Context, username string, password string) (int, error) {
	const op = "services.auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op))
	log.Info("registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := a.userSaver.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user is already exists", err.Error())
			return 0, fmt.Errorf("%s %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to save user", err.Error())
		return 0, fmt.Errorf("%s %w", op, err)
	}
	return id, nil
}

func (a *AuthService) Authorize(tokenString string) (string, error) {
	const op = "services.auth.Authorize"
	log := a.log.With(slog.String("op", op))
	log.Info("authorize user")
	username, err := jwt.ParseToken(tokenString, a.secret)
	if err != nil {
		log.Error("failed to authorize user", err.Error())
		return "", fmt.Errorf("%s: %w")
	}
	return username, nil
}

func (a *AuthService) GetCoinsInfo(ctx context.Context, username string) (int, error) {
	const op = "services.auth.GetCoinsInfo"
	log := a.log.With(slog.String("op", op))
	log.Info("getting coins info")
	user, err := a.userProvider.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return 0, fmt.Errorf("%s %w", op, ErrInvalidCredentials)
		}
		return 0, fmt.Errorf("%s %w", op, err)
	}
	return user.Coins, nil
}
