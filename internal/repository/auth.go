package repository

import (
	"context"
	"time"

	"example.com/mod/internal/domain"
)

type AuthRepository interface {
	SaveToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*domain.RefreshToken, error)
	GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	DeleteByHash(ctx context.Context, tokenHash string) error
}
