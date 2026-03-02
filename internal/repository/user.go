package repository

import (
	"context"
	"example.com/mod/internal/domain"
)

type UserRepository interface {
	GetByEmail (ctx context.Context, email string)(*domain.User, error)
	GetByID (ctx context.Context, id int) (*domain.User, error)
	Create (ctx context.Context, email, passwordHash, name, role string) (*domain.User, error)
}
