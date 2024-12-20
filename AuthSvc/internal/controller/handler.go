package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/storage"
	"net/http"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User) (int, error)
	LoginUser(ctx context.Context, user *models.User) (string, error)
	IsHotelier(ctx context.Context, userID int) (bool, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	_, err := a.authService.RegisterUser(ctx, &user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//TODO: ErrInvalidCredentials

	token, err := a.authService.LoginUser(ctx, &user)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
