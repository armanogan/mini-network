package token

import (
	"mini-network/internal/v1/models"
	"time"
)

type TokenService interface {
	ValidateToken(tokenString string) error
	GenerateToken(user models.User) (string, error)
	ParseToken(tokenString string) (models.User, error)
	GetExpired() time.Duration
}
