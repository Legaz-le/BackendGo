package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"example.com/mod/internal/auth"
	"example.com/mod/internal/domain"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "testsecret")
	os.Exit(m.Run())
}

// --- mocks ---

type mockUserRepo struct {
	user *domain.User
	err  error
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) Create(ctx context.Context, email, passwordHash, name, role string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.User{ID: 1, Email: email, Name: name, Role: role}, nil
}

type mockAuthRepo struct {
	token *domain.RefreshToken
	err   error
}

func (m *mockAuthRepo) SaveToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) (*domain.RefreshToken, error) {
	return m.token, m.err
}

func (m *mockAuthRepo) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	return m.token, m.err
}

func (m *mockAuthRepo) DeleteByHash(ctx context.Context, tokenHash string) error {
	return m.err
}

// --- Register ---

func TestAuthService_Register(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{}, &mockAuthRepo{})

	user, accessToken, rawToken, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User", "seeker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", user.Email)
	}
	if accessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if rawToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
}

func TestAuthService_Register_UserCreateError(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{err: errors.New("duplicate email")}, &mockAuthRepo{})

	_, _, _, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User", "seeker")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuthService_Register_SaveTokenError(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{}, &mockAuthRepo{err: errors.New("db error")})

	_, _, _, err := svc.Register(context.Background(), "test@example.com", "password123", "Test User", "seeker")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Login ---

func TestAuthService_Login(t *testing.T) {
	hash, _ := auth.HashPassword("password123")
	user := &domain.User{ID: 1, Email: "test@example.com", Role: "seeker", PasswordHash: hash}
	svc := NewAuthService(&mockUserRepo{user: user}, &mockAuthRepo{})

	gotUser, accessToken, rawToken, err := svc.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUser.ID != 1 {
		t.Fatalf("expected user ID 1, got %d", gotUser.ID)
	}
	if accessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if rawToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{err: errors.New("not found")}, &mockAuthRepo{})

	_, _, _, err := svc.Login(context.Background(), "missing@example.com", "password123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hash, _ := auth.HashPassword("password123")
	user := &domain.User{ID: 1, Email: "test@example.com", Role: "seeker", PasswordHash: hash}
	svc := NewAuthService(&mockUserRepo{user: user}, &mockAuthRepo{})

	_, _, _, err := svc.Login(context.Background(), "test@example.com", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

// --- Refresh ---

func TestAuthService_Refresh(t *testing.T) {
	hash, _ := auth.HashPassword("password123")
	user := &domain.User{ID: 1, Email: "test@example.com", Role: "seeker", PasswordHash: hash}
	token := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: "somehash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	svc := NewAuthService(&mockUserRepo{user: user}, &mockAuthRepo{token: token})

	gotUser, accessToken, rawToken, err := svc.Refresh(context.Background(), "somehash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUser.ID != 1 {
		t.Fatalf("expected user ID 1, got %d", gotUser.ID)
	}
	if accessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if rawToken == "" {
		t.Fatal("expected non-empty refresh token")
	}
}

func TestAuthService_Refresh_TokenNotFound(t *testing.T) {
	svc := NewAuthService(&mockUserRepo{}, &mockAuthRepo{err: errors.New("token not found")})

	_, _, _, err := svc.Refresh(context.Background(), "badhash")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuthService_Refresh_TokenExpired(t *testing.T) {
	token := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		TokenHash: "somehash",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	svc := NewAuthService(&mockUserRepo{}, &mockAuthRepo{token: token})

	_, _, _, err := svc.Refresh(context.Background(), "somehash")
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	if err.Error() != "token expired" {
		t.Fatalf("expected 'token expired', got '%s'", err.Error())
	}
}
