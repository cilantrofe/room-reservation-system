package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/models"
)

type Storage interface {
	RegisterUser(ctx context.Context, user *models.User) (int, error)
	LoginUser(ctx context.Context, chatID string) (*models.User, error)
	IsHotelier(ctx context.Context, userID int) (bool, error)
}
