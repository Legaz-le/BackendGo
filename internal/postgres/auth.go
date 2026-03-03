package postgres

import (
	"context"
	"errors"
	"time"

	"example.com/mod/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuthRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAuthRepository(db *pgxpool.Pool) *PostgresAuthRepository {
	return &PostgresAuthRepository{db: db}
}

var ErrTokenNotFound = errors.New("token not found")

func (r *PostgresAuthRepository) SaveToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*domain.RefreshToken, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token_hash, expires_at", userID, tokenHash, expiresAt)
	token := domain.RefreshToken{}

	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *PostgresAuthRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1", tokenHash)
	token := domain.RefreshToken{}

	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *PostgresAuthRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	pool := r.db

	result, err := pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash = $1", tokenHash)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTokenNotFound
	}

	return nil
}
