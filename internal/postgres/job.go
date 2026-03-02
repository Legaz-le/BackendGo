package postgres

import (
	"context"
	"errors"
	"strconv"

	"example.com/mod/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresJobRepository struct {
	db *pgxpool.Pool
}

func NewPostgresJobRepository(db *pgxpool.Pool) *PostgresJobRepository {
	return &PostgresJobRepository{db: db}
}

var ErrJobNotFound = errors.New("job not found")

func (r *PostgresJobRepository) GetAll(ctx context.Context) ([]domain.Job, error) {
	pool := r.db

	rows, err := pool.Query(ctx, "SELECT id, title, description, location, salary_min, salary_max, created_at FROM jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []domain.Job

	for rows.Next() {
		job := domain.Job{}
		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.SalaryMin, &job.SalaryMax, &job.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *PostgresJobRepository) GetByID(ctx context.Context, id int) (*domain.Job, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "SELECT id, title, description, location, salary_min, salary_max, created_at FROM jobs WHERE id = $1", id)

	job := domain.Job{}
	err := row.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.SalaryMin, &job.SalaryMax, &job.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return &job, nil

}

func (r *PostgresJobRepository) Create(ctx context.Context, job domain.Job) (*domain.Job, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "INSERT INTO jobs (title, description, location, salary_min, salary_max) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at",
		job.Title, job.Description, job.Location, job.SalaryMin, job.SalaryMax)

	err := row.Scan(&job.ID, &job.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *PostgresJobRepository) Update(ctx context.Context, id int, job domain.Job) (*domain.Job, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "UPDATE jobs SET title = $1, description = $2, location = $3, salary_min = $4, salary_max = $5 WHERE id = $6 RETURNING id, title, description, location, salary_min, salary_max, created_at",
		job.Title, job.Description, job.Location, job.SalaryMin, job.SalaryMax, id)

	updatedJob := domain.Job{}
	err := row.Scan(&updatedJob.ID, &updatedJob.Title, &updatedJob.Description, &updatedJob.Location, &updatedJob.SalaryMin, &updatedJob.SalaryMax, &updatedJob.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return &updatedJob, nil
}

func (r *PostgresJobRepository) Delete(ctx context.Context, id int) error {
	pool := r.db

	result, err := pool.Exec(ctx, "DELETE FROM jobs WHERE id = $1", id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrJobNotFound
	}

	return nil
}

func (r *PostgresJobRepository) GetWithFilter(ctx context.Context, location string, minSalary int, maxSalary int) ([]domain.Job, error) {
	pool := r.db

	query := "SELECT id, title, description, location, salary_min, salary_max, created_at FROM jobs WHERE 1=1"
	var args []interface{}
	paramCount := 1

	if location != "" {
		query += " AND location ILIKE $" + strconv.Itoa(paramCount)
		args = append(args, "%"+location+"%")
		paramCount++
	}

	if minSalary > 0 {
		query += " AND salary_min >= $" + strconv.Itoa(paramCount)
		args = append(args, minSalary)
		paramCount++
	}

	if maxSalary > 0 {
		query += " AND salary_max <= $" + strconv.Itoa(paramCount)
		args = append(args, maxSalary)
		paramCount++
	}

	result, err := pool.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	var jobs []domain.Job
	for result.Next() {
		job := domain.Job{}
		if err := result.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.SalaryMin, &job.SalaryMax, &job.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return jobs, nil

}
