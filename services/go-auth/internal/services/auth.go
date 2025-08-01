package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-auth/internal/app"
	"go-auth/internal/models"
	"go-auth/internal/storage"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userStorage app.AppUserStorage
	tokenStore  app.AppTokenStorage
}

func AuthNew(userStorage app.AppUserStorage, tokenStore app.AppTokenStorage) *AuthService {
	return &AuthService{userStorage: userStorage, tokenStore: tokenStore}
}

func (s AuthService) Create(user models.UserCreateReq) (*models.TokenDto, error) {
	jwt := models.TokenDto{}
	if user.Login == "" || user.Password == "" {
		return nil, errors.New("login and password are required")
	}
	pswdHash, err := getHash(user.Password)
	if err != nil {
		return nil, err
	}
	existingUser, err := s.userStorage.Get(models.UserCreateDto{Login: user.Login, PasswordHash: pswdHash})

	if err != nil && err != storage.ErrUserNotFound {
		return nil, fmt.Errorf("AuthService Create check existingUser: %v", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}
	u, err := s.userStorage.Save(models.UserCreateDto{Login: user.Login, PasswordHash: pswdHash})
	if err != nil {
		return nil, fmt.Errorf("AuthService Create userStorage Save: %v", err)
	}
	// TODO refresh and access tokens

	access, err := s.generateJWT(u)
	if err != nil {
		return nil, err
	}
	refresh, err := s.generateJWT(u)
	if err != nil {
		return nil, err
	}
	jwt.Access = access
	jwt.Refresh = refresh

	return &jwt, nil
}

func (s AuthService) Login(user models.UserCreateReq) (*models.TokenDto, error) {
	jwt := models.TokenDto{}
	if user.Login == "" || user.Password == "" {
		return nil, errors.New("login and password are required")
	}
	pswdHash, err := getHash(user.Password)
	if err != nil {
		return nil, err
	}
	existingUser, err := s.userStorage.Get(models.UserCreateDto{Login: user.Login, PasswordHash: pswdHash})
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, ErrWrongLoginOrPassword
	}
	// TODO refresh and access tokens

	access, err := s.generateJWT(*existingUser)
	if err != nil {
		return nil, err
	}
	refresh, err := s.generateJWT(*existingUser)
	if err != nil {
		return nil, err
	}
	jwt.Access = access
	jwt.Refresh = refresh

	return &jwt, nil
}

func (s AuthService) RefreshToken(refreshToken string) (*models.TokenDto, error) {
	// 1. Валидация refresh token
	var secret string
	err := app.AppContainer.Invoke(func(cfg app.AppConfig) {
		secret = cfg.GetConfig().Secret
	})
	if err != nil {
		return nil, fmt.Errorf("error getting config: %w", err)
	}

	token, err := validateJWT(refreshToken, secret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 2. Проверка срока действия
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("refresh token expired")
	}

	// 3. Получение информации о пользователе
	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid subject claim")
	}

	user, err := s.userStorage.GetById(sub)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 4. Генерация новых токенов
	newAccess, err := s.generateJWT(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := s.generateJWT(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// 5. Обновление токенов в хранилище
	newTokens := &models.TokenDto{
		Access:  newAccess,
		Refresh: newRefresh,
	}

	// Если есть tokenStore, сохраняем новые токены
	if s.tokenStore != nil {
		ctx := context.Background()
		err = s.tokenStore.SetTokens(ctx, newTokens)
		if err != nil {
			return nil, fmt.Errorf("failed to store tokens: %w", err)
		}

		// Удаляем старый refresh token
		err = s.tokenStore.RemoveToken(ctx, refreshToken, "")
		if err != nil {
			return nil, fmt.Errorf("failed to remove old token: %w", err)
		}
	}

	return newTokens, nil
}
func (s AuthService) generateJWT(user models.UserCreateRes) (string, error) {
	var secret string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})
	err := app.AppContainer.Invoke(func(cfg app.AppConfig) {
		secret = cfg.GetConfig().Secret
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
