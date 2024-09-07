package usecase

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"mini-network/internal/v1/models"
	"mini-network/internal/v1/repository"
	"mini-network/internal/v1/token"
	errors2 "mini-network/pkg/errors"
	"net/http"
)

type User interface {
	GetUserByID(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, name, email, password string) (int, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	GetTokenService() token.TokenService
}

type userUseCase struct {
	tokenService token.TokenService
	repo         repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository, tokenService token.TokenService) *userUseCase {
	return &userUseCase{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (u *userUseCase) GetTokenService() token.TokenService {
	return u.tokenService
}

func (u *userUseCase) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *userUseCase) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}
func (u *userUseCase) CreateUser(ctx context.Context, name, email, password string) (int, error) {
	user, err := u.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, errors2.NewErrorWithCode(http.StatusInternalServerError, err.Error())
	}
	if user.ID > 0 {
		return 0, errors2.NewErrorWithCode(http.StatusConflict, "Email exist")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors2.NewErrorWithCode(http.StatusInternalServerError, err.Error())
	}
	return u.repo.SaveUser(ctx, email, name, string(hash))
}

func (u *userUseCase) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := u.GetUserByEmail(ctx, email)
	if user.ID == 0 {
		return "", errors2.NewErrorWithCode(http.StatusUnauthorized, err.Error())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors2.NewErrorWithCode(http.StatusUnauthorized, err.Error())
	}
	token, err := u.tokenService.GenerateToken(user)
	if err != nil {
		return "", errors2.NewErrorWithCode(http.StatusInternalServerError, err.Error())
	}
	return token, err
}
