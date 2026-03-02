package service

import (
	"context"

	"example.com/mod/internal/domain"
	"example.com/mod/internal/repository"
)

type JobService struct {
	repo repository.JobRepository
}

func NewJobService(repo repository.JobRepository) *JobService {
	return &JobService{repo: repo}
}

func (s *JobService) GetAll(ctx context.Context) ([]domain.Job, error) {
	return s.repo.GetAll(ctx)
}

func (s *JobService) GetByID(ctx context.Context, id int) (*domain.Job, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *JobService) Create(ctx context.Context, job domain.Job) (*domain.Job, error) {
	return s.repo.Create(ctx, job)
}
func (s *JobService) Update(ctx context.Context, id int, job domain.Job) (*domain.Job, error) {
	return s.repo.Update(ctx, id, job)
}
func (s *JobService) Delete(ctx context.Context, id int)  error {
	return s.repo.Delete(ctx, id)
	
}
func (s *JobService) GetWithFilter(ctx context.Context, location string, minSalary int, maxSalary int) ([]domain.Job, error) {
	return s.repo.GetWithFilter(ctx, location, minSalary, maxSalary)
}
