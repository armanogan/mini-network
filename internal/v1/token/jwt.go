package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"mini-network/internal/v1/models"
	"time"
)

type Jwt struct {
	secret  string
	expired time.Duration
}

func NewJwt(secret string, expired time.Duration) *Jwt {
	return &Jwt{secret: secret, expired: expired}
}
func (j Jwt) GetExpired() time.Duration {
	return j.expired
}
func (j Jwt) GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(j.expired).Unix(), // Токен истекает через 72 часа
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j Jwt) ValidateToken(tokenString string) error {
	_, err := j.parse(tokenString)
	return err
}
func (j Jwt) parse(tokenString string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (j Jwt) ParseToken(tokenString string) (models.User, error) {
	var user models.User
	claims, err := j.parse(tokenString)
	if err != nil {
		return user, err
	}
	userId, ok := claims["sub"].(float64)
	if !ok {
		return user, fmt.Errorf("Error converting user_id")
	}
	name, ok := claims["name"].(string)
	if !ok {
		return user, fmt.Errorf("Error converting name")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return user, fmt.Errorf("Error converting email")
	}
	user.ID = int(userId)
	user.Name = name
	user.Email = email
	return user, nil
}
