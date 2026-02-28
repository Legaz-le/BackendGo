package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/mod/internal/user"
	"github.com/go-playground/validator/v10"
)

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required" `
	Role     string `json:"role" validate:"required,oneof=employer seeker"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	hash, err := HashPassword(req.Password)

	if err != nil {
		serverError(w, err)
		return
	}

	createUser, err := user.CreateUser(r.Context(), req.Email, hash, req.Name, req.Role)

	if err != nil {
		serverError(w, err)
		return
	}

	accessToken, err := GenerateAccessToken(createUser.ID, createUser.Role)

	if err != nil {
		serverError(w, err)
		return
	}

	rawToken, tokenHash, err := GenerateRefreshToken()

	if err != nil {
		serverError(w, err)
		return
	}

	_, err = SaveToken(r.Context(), createUser.ID, tokenHash, time.Now().Add(7*24*time.Hour))

	if err != nil {
		serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rawToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})
	w.WriteHeader(http.StatusOK)

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	foundUser, err := user.GetUserByEmail(r.Context(), req.Email)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	err = CheckPasswordHash(req.Password, foundUser.PasswordHash)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := GenerateAccessToken(foundUser.ID, foundUser.Role)

	if err != nil {
		serverError(w, err)
		return
	}

	rawToken, tokenHash, err := GenerateRefreshToken()

	if err != nil {
		serverError(w, err)
		return
	}

	_, err = SaveToken(r.Context(), foundUser.ID, tokenHash, time.Now().Add(7*24*time.Hour))

	if err != nil {
		serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rawToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})
	w.WriteHeader(http.StatusOK)

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var rawToken string

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tokenHash := HashRefreshToken(cookie.Value)

	getTokenHash, err := GetRefreshTokenByHash(r.Context(), tokenHash)

	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	if getTokenHash.ExpiresAt.Before(time.Now()) {
		http.Error(w, "token expired", http.StatusUnauthorized)
		return
	}

	err = DeleteRefreshToken(r.Context(), tokenHash)
	if err != nil {
		serverError(w, err)
		return
	}

	foundUser, err := user.GetUserByID(r.Context(), getTokenHash.UserID)

	if err != nil {
		serverError(w, err)
		return
	}

	accessToken, err := GenerateAccessToken(foundUser.ID, foundUser.Role)

	if err != nil {
		serverError(w, err)
		return
	}

	rawToken, tokenHash, err = GenerateRefreshToken()

	if err != nil {
		serverError(w, err)
		return
	}

	_, err = SaveToken(r.Context(), foundUser.ID, tokenHash, time.Now().Add(7*24*time.Hour))
	if err != nil {
		serverError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    rawToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})
	w.WriteHeader(http.StatusOK)
}

func Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	claims, err := ValidateAccessToken(cookie.Value)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "refresh_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	w.WriteHeader(http.StatusOK)

}
