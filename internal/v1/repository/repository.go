package repository

import (
	"context"
	"mini-network/internal/v1/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	SaveUser(ctx context.Context, email, name, password string) (int, error)
}
