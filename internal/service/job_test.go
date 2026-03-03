package service

import (
	"context"
	"errors"
	"testing"

	"example.com/mod/internal/domain"
)

type mockJobRepo struct {
	jobs []domain.Job
	err  error
}

func (m *mockJobRepo) GetAll(ctx context.Context) ([]domain.Job, error) {
	return m.jobs, m.err
}

func (m *mockJobRepo) GetByID(ctx context.Context, id int) (*domain.Job, error) {
	for _, j := range m.jobs {
		if j.ID == id {
			return &j, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockJobRepo) Create(ctx context.Context, job domain.Job) (*domain.Job, error) {
	if m.err != nil {
		return nil, m.err
	}
	job.ID = 1
	return &job, nil
}

func (m *mockJobRepo) Update(ctx context.Context, id int, job domain.Job) (*domain.Job, error) {
	if m.err != nil {
		return nil, m.err
	}
	job.ID = id
	return &job, nil
}

func (m *mockJobRepo) Delete(ctx context.Context, id int) error {
	return m.err
}

func (m *mockJobRepo) GetWithFilter(ctx context.Context, location string, minSalary int, maxSalary int) ([]domain.Job, error) {
	return m.jobs, m.err
}

func TestJobService_GetAll(t *testing.T) {
	jobs := []domain.Job{
		{ID: 1, Title: "Go Engineer", Location: "Remote"},
		{ID: 2, Title: "Frontend Dev", Location: "NYC"},
	}
	svc := NewJobService(&mockJobRepo{jobs: jobs})

	result, err := svc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(result))
	}
}

func TestJobService_GetAll_Error(t *testing.T) {
	svc := NewJobService(&mockJobRepo{err: errors.New("db error")})

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJobService_GetByID(t *testing.T) {
	jobs := []domain.Job{{ID: 1, Title: "Go Engineer", Location: "Remote"}}
	svc := NewJobService(&mockJobRepo{jobs: jobs})

	job, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID != 1 {
		t.Fatalf("expected ID 1, got %d", job.ID)
	}
}

func TestJobService_GetByID_NotFound(t *testing.T) {
	svc := NewJobService(&mockJobRepo{jobs: []domain.Job{}})

	_, err := svc.GetByID(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJobService_Create(t *testing.T) {
	svc := NewJobService(&mockJobRepo{})

	job, err := svc.Create(context.Background(), domain.Job{
		Title:       "Go Engineer",
		Description: "Build APIs",
		Location:    "Remote",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID != 1 {
		t.Fatalf("expected ID 1, got %d", job.ID)
	}
	if job.Title != "Go Engineer" {
		t.Fatalf("expected title 'Go Engineer', got %s", job.Title)
	}
}

func TestJobService_Create_Error(t *testing.T) {
	svc := NewJobService(&mockJobRepo{err: errors.New("db error")})

	_, err := svc.Create(context.Background(), domain.Job{Title: "Go Engineer"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJobService_Update(t *testing.T) {
	svc := NewJobService(&mockJobRepo{})

	job, err := svc.Update(context.Background(), 1, domain.Job{
		Title:       "Updated Title",
		Description: "Updated Desc",
		Location:    "NYC",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if job.ID != 1 {
		t.Fatalf("expected ID 1, got %d", job.ID)
	}
	if job.Title != "Updated Title" {
		t.Fatalf("expected updated title, got %s", job.Title)
	}
}

func TestJobService_Update_Error(t *testing.T) {
	svc := NewJobService(&mockJobRepo{err: errors.New("not found")})

	_, err := svc.Update(context.Background(), 99, domain.Job{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJobService_Delete(t *testing.T) {
	svc := NewJobService(&mockJobRepo{})

	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJobService_Delete_Error(t *testing.T) {
	svc := NewJobService(&mockJobRepo{err: errors.New("not found")})

	err := svc.Delete(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJobService_GetWithFilter(t *testing.T) {
	jobs := []domain.Job{
		{ID: 1, Title: "Go Engineer", Location: "Remote", SalaryMin: 80000, SalaryMax: 120000},
	}
	svc := NewJobService(&mockJobRepo{jobs: jobs})

	result, err := svc.GetWithFilter(context.Background(), "Remote", 70000, 130000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 job, got %d", len(result))
	}
}
