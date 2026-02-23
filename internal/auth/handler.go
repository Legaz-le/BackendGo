package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/mod/internal/user"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type authResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse{AccessToken: accessToken, RefreshToken: rawToken})

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{AccessToken: accessToken, RefreshToken: rawToken})

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	var rawToken string

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tokenHash := HashRefreshToken(req.RefreshToken)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{AccessToken: accessToken, RefreshToken: rawToken})
}
