package auth

import (
	"context"
	"errors"
	"time"

	"example.com/mod/internal/database"
	"example.com/mod/internal/models"
)

var ErrTokenNotFound = errors.New("token not found")

func SaveToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*models.RefreshToken, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token_hash, expires_at", userID, tokenHash, expiresAt)
	token := models.RefreshToken{}

	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1", tokenHash)
	token := models.RefreshToken{}

	err := row.Scan(&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &token, nil
}

func DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	pool := database.GetDB()

	result, err := pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash = $1", tokenHash)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTokenNotFound
	}
	
	return nil
}
