package services

import (
	"errors"
	"fmt"
	"time"

	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
}

func AuthNew() *AuthService {
	return &AuthService{}
}

func (s AuthService) Create(user models.UserCreate) (string, error) {
	u := models.UserCreateRes{}
	if user.Login == "" || user.Password == "" {
		return "", errors.New("login and password are required")
	}
	// если в БД их нет, то записать

	jwt, err := s.GenerateJWT(u)
	if err != nil {
		return jwt, err
	}
	return jwt, nil
}
func (s AuthService) GenerateJWT(user models.UserCreateRes) (string, error) {
	var secret string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})
	err := app.AppContainer.Invoke(func(cfg *config.Config) {
		secret = cfg.Secret
	})
	if err != nil {
		return "", fmt.Errorf("error with invoke config, %w", err)
	}
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
