package services

import (
	"errors"
	"fmt"
	"time"

	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/models"
	"go-auth/internal/storage"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userStorage *storage.UserStorage
}

func AuthNew(userStorage *storage.UserStorage) *AuthService {
	return &AuthService{userStorage: userStorage}
}

func (s AuthService) Create(user models.UserCreateReq) (string, error) {

	if user.Login == "" || user.Password == "" {
		return "", errors.New("login and password are required")
	}
	pswdHash, err := getHash(user.Password)
	if err != nil {
		return "", err
	}
	existingUser, err := s.userStorage.Get(models.UserCreateDto{Login: user.Login, PasswordHash: pswdHash})
	if err != nil {
		return "", err
	}
	if existingUser != nil {
		return "", ErrUserExists
	}
	u := s.userStorage.Save(user)
	jwt, err := s.generateJWT(u)
	if err != nil {
		return jwt, err
	}
	return jwt, nil
}

func (s AuthService) generateJWT(user models.UserCreateRes) (string, error) {
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

func validateJWT(tokenString string, secret string) (*jwt.Token, error) {
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

func getHash(s string) (string, error) {
	return "", nil
}
