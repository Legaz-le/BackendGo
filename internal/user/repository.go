package user

import (
	"context"

	"example.com/mod/internal/database"
	"example.com/mod/internal/models"
)

func CreateUser(ctx context.Context, email, passwordHash, name, role string) (*models.User, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "INSERT INTO users (email, password_hash, name, role) VALUES ($1, $2, $3, $4) RETURNING id, password_hash, email, name, role, created_at",
		email, passwordHash, name, role)
	user := models.User{}
	err := row.Scan(&user.ID, &user.PasswordHash, &user.Email, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "SELECT id, email, password_hash, name, role, created_at FROM users WHERE email = $1", email)
	user := models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	pool := database.GetDB()

	row := pool.QueryRow(ctx, "SELECT id, email, password_hash, name, role, created_at FROM users WHERE id = $1", id)
	user := models.User{}

	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
