package postgres

import (
	"context"

	"example.com/mod/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}



func (r *PostgresUserRepository) Create(ctx context.Context, email, passwordHash, name, role string) (*domain.User, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, $4) RETURNING id, password_hash, email, name, role, created_at",
		email, passwordHash, name, role)
	user := domain.User{}
	err := row.Scan(&user.ID, &user.PasswordHash, &user.Email, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "SELECT id, email, password_hash, name, role, created_at FROM users WHERE email = $1", email)
	user := domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	pool := r.db

	row := pool.QueryRow(ctx, "SELECT id, email, password_hash, name, role, created_at FROM users WHERE id = $1", id)
	user := domain.User{}

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
