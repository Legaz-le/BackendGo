package models

import (
	"context"
	"errors"
	"strconv"

	"example.com/mod/internal/database"
	"github.com/jackc/pgx/v5"
)

var ErrJobNotFound = errors.New("job not found")

func GetAllJobs(ctx context.Context) ([]Job, error) {
	pool := database.GetDB()

	rows, err := pool.Query(ctx, "SELECT id, title, description, location, salary_range, created_at FROM jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job

	for rows.Next() {
		job := Job{}
		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.SalaryRange, &job.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetJobByID(ctx context.Context, id int) (*Job, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "SELECT id, title, description, location, salary_range, created_at FROM jobs WHERE id = $1", id)

	job := Job{}
	err := row.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.SalaryRange, &job.CreatedAt)

	if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrJobNotFound
			}
			return nil, err
		}

	return &job, nil

}

func CreateJob(ctx context.Context, job Job) (*Job, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "INSERT INTO jobs (title, description, location, salary_range) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		job.Title, job.Description, job.Location, job.SalaryRange)

	err := row.Scan(&job.ID, &job.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func UpdateJob(ctx context.Context, id int, job Job) (*Job, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "UPDATE jobs SET title = $1, description = $2, location = $3, salary_range = $4 WHERE id = $5 RETURNING id, title, description, location, salary_range, created_at",
		job.Title, job.Description, job.Location, job.SalaryRange, id)

	updatedJob := Job{}
	err := row.Scan(&updatedJob.ID, &updatedJob.Title, &updatedJob.Description, &updatedJob.Location, &updatedJob.SalaryRange, &updatedJob.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}
		return nil, err
	}

	return &updatedJob, nil
}

func DeleteJob(ctx context.Context, id int) error {
	pool := database.GetDB()

	result, err := pool.Exec(ctx, "DELETE FROM jobs WHERE id = $1", id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrJobNotFound
	}

	return nil
}

func GetJobsWithFilter(ctx context.Context, location string, minSalary int, maxSalary int) ([]Job, error) {
	pool := database.GetDB()

	query := "SELECT id, title, description, location, salary_range, created_at FROM jobs WHERE 1=1"
	var args []interface{}
	paramCount := 1

	if location != "" {
		query += " AND location ILIKE $" + strconv.Itoa(paramCount)
		args = append(args, "%"+location+"%")
		paramCount++
	}

	if minSalary > 0 {
		query += " AND salary_range::numeric >= $" + strconv.Itoa(paramCount)
		args = append(args, minSalary)
		paramCount++
	}

	if maxSalary > 0 {
		query += " AND salary_range::numeric <= $" + strconv.Itoa(paramCount)
		args = append(args, maxSalary)
		paramCount++
	}

	result, err := pool.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	var jobs []Job
	for result.Next() {
		job := Job{}
		if err := result.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.SalaryRange, &job.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return jobs, nil

}
