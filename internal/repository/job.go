package repository

import (
	"context"
	"example.com/mod/internal/domain"
)

type JobRepository interface {
	GetAll(ctx context.Context) ([]domain.Job, error)
	GetByID(ctx context.Context, id int) (*domain.Job, error)
	Create(ctx context.Context, job domain.Job) (*domain.Job, error)
	Update(ctx context.Context, id int, job domain.Job) (*domain.Job, error)
	Delete(ctx context.Context, id int) error
	GetWithFilter(ctx context.Context, location string, minSalary int, maxSalary int) ([]domain.Job, error)
}
