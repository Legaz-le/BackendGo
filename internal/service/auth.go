package service

import (
	"context"
	"errors"
	"time"

	"example.com/mod/internal/auth"
	"example.com/mod/internal/domain"
	"example.com/mod/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewAuthService(userRepo repository.UserRepository, authRepo repository.AuthRepository) *AuthService {
	return &AuthService{userRepo: userRepo, authRepo: authRepo}
}

func (r *AuthService) Register(ctx context.Context, email, password, name, role string) (*domain.User, string, string, error) {
	hash, err := auth.HashPassword(password)

	if err != nil {
		return nil, "", "", err
	}
	user, err := r.userRepo.Create(ctx, email, hash, name, role)

	if err != nil {
		return nil, "", "", err
	}
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)

	if err != nil {
		return nil, "", "", err
	}
	rawToken, tokenHash, err := auth.GenerateRefreshToken()

	if err != nil {
		return nil, "", "", err
	}
	_, err = r.authRepo.SaveToken(ctx, user.ID, tokenHash, time.Now().Add(7*24*time.Hour))

	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, rawToken, nil
}

func (r *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, string, error) {
	user, err := r.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return nil, "", "", err
	}

	checkPassword := auth.CheckPasswordHash(password, user.PasswordHash)

	if checkPassword != nil {
		return nil, "", "", checkPassword
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)

	if err != nil {
		return nil, "", "", err
	}

	rawToken, tokenHash, err := auth.GenerateRefreshToken()

	if err != nil {
		return nil, "", "", err
	}

	_, err = r.authRepo.SaveToken(ctx, user.ID, tokenHash, time.Now().Add(7*24*time.Hour))

	if err != nil {
		return nil, "", "", err

	}
	return user, accessToken, rawToken, nil
}

func (r *AuthService) Refresh(ctx context.Context, tokenHash string) (*domain.User, string, string, error) {
	token, err := r.authRepo.GetByHash(ctx, tokenHash)

	if err != nil {
		return nil, "", "", err
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, "", "", errors.New("token expired")
	}

	err = r.authRepo.DeleteByHash(ctx, tokenHash)
	
	if err != nil {
		return nil, "", "", err
	}
	
	user, err := r.userRepo.GetByID(ctx, token.UserID)
	
	if err != nil {
		return nil, "", "", err
	}
	
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)

	if err != nil {
		return nil, "", "", err
	}

	rawToken, tokenHash, err := auth.GenerateRefreshToken()

	if err != nil {
		return nil, "", "", err
	}

	_, err = r.authRepo.SaveToken(ctx, user.ID, tokenHash, time.Now().Add(7*24*time.Hour))

	if err != nil {
		return nil, "", "", err

	}

	return user, accessToken, rawToken, nil
}
