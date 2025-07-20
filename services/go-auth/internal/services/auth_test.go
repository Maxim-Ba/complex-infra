package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/dig"
)

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) GetById(userId string) (*models.UserCreateRes, error) {
	args := m.Called(userId)
	return args.Get(0).(*models.UserCreateRes), args.Error(1)
}

func (m *MockUserStorage) Update(user models.UserCreateDto) error {
	args := m.Called(user)
	return args.Error(1)
}

func (m *MockUserStorage) Save(user models.UserCreateDto) (models.UserCreateRes, error) {
	args := m.Called(user)
	return args.Get(0).(models.UserCreateRes), args.Error(1)
}

func (m *MockUserStorage) Get(user models.UserCreateDto) (*models.UserCreateRes, error) {
	args := m.Called(user)
	return args.Get(0).(*models.UserCreateRes), args.Error(1)
}

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetConfig() *config.Config {
	args := m.Called()
	return args.Get(0).(*config.Config)
}

type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) SetTokens(ctx context.Context, tokens *models.TokenDto) error {
	args := m.Called(ctx, tokens)
	return args.Error(0)
}

func (m *MockTokenStorage) RemoveToken(ctx context.Context, refreshToken string, accessToken string) error {
	args := m.Called(ctx, refreshToken, accessToken)
	return args.Error(0)
}

func (m *MockTokenStorage) GetTokens(ctx context.Context, refreshToken string) (*models.TokenDto, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*models.TokenDto), args.Error(1)
}

func TestAuthService_Create(t *testing.T) {
	originalContainer := app.AppContainer
	defer func() { app.AppContainer = originalContainer }()

	testContainer := dig.New()
	app.AppContainer = testContainer

	tests := []struct {
		name          string
		input         models.UserCreateReq
		mockSetup     func(*MockUserStorage, *MockConfig)
		expectedError error
		wantToken     bool
	}{
		{
			name: "successful user creation",
			input: models.UserCreateReq{
				Login:    "testuser",
				Password: "password123",
			},
			mockSetup: func(mus *MockUserStorage, mc *MockConfig) {
				mus.On("Get", mock.Anything).Return(
					(*models.UserCreateRes)(nil),
					nil,
				)

				mus.On("Save", mock.Anything).Return(
					models.UserCreateRes{
						Login: "testuser",
						Id:    "123",
					}, nil,
				)

				mc.On("GetConfig").Return(
					&config.Config{
						Secret: "test-secret",
					},
				)
			},
			expectedError: nil,
			wantToken:     true,
		},
		{
			name: "empty login and password",
			input: models.UserCreateReq{
				Login:    "",
				Password: "",
			},
			mockSetup:     func(mus *MockUserStorage, mc *MockConfig) {},
			expectedError: errors.New("login and password are required"),
			wantToken:     false,
		},
		{
			name: "user already exists",
			input: models.UserCreateReq{
				Login:    "existing",
				Password: "password",
			},
			mockSetup: func(mus *MockUserStorage, mc *MockConfig) {
				mus.On("Get", mock.Anything).Return(
					&models.UserCreateRes{
						Login: "existing",
						Id:    "456",
					},
					nil,
				)
			},
			expectedError: ErrUserExists,
			wantToken:     false,
		},
		{
			name: "error getting user",
			input: models.UserCreateReq{
				Login:    "erroruser",
				Password: "password",
			},
			mockSetup: func(mus *MockUserStorage, mc *MockConfig) {
				mus.On("Get", mock.Anything).Return(
					(*models.UserCreateRes)(nil),
					errors.New("database error"),
				)
			},
			expectedError: errors.New("AuthService Create check existingUser: database error"),
			wantToken:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStorage := &MockUserStorage{}
			mockConfig := &MockConfig{}
			mockTokenStorage := &MockTokenStorage{}

			tt.mockSetup(mockUserStorage, mockConfig)

			// Очищаем контейнер перед каждым тестом
			testContainer = dig.New()
			app.AppContainer = testContainer

			err := testContainer.Provide(func() app.AppConfig {
				return mockConfig
			})
			assert.NoError(t, err)

			service := AuthNew(mockUserStorage, mockTokenStorage)

			token, err := service.Create(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotNil(t, token)
					assert.NotEmpty(t, token.Access)
					assert.NotEmpty(t, token.Refresh)
				} else {
					assert.Nil(t, token)
				}
			}

			mockUserStorage.AssertExpectations(t)
			mockConfig.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	originalContainer := app.AppContainer
	defer func() { app.AppContainer = originalContainer }()

	testContainer := dig.New()
	app.AppContainer = testContainer

	tests := []struct {
		name          string
		input         models.UserCreateReq
		mockSetup     func(*MockUserStorage, *MockConfig)
		expectedError error
		wantToken     bool
	}{
		{
			name: "successful login",
			input: models.UserCreateReq{
				Login:    "testuser",
				Password: "password123",
			},
			mockSetup: func(mus *MockUserStorage, mc *MockConfig) {
				pswdHash, _ := getHash("password123")
				mus.On("Get", models.UserCreateDto{
					Login:        "testuser",
					PasswordHash: pswdHash,
				}).Return(
					&models.UserCreateRes{
						Login: "testuser",
						Id:    "123",
					},
					nil,
				)

				mc.On("GetConfig").Return(
					&config.Config{
						Secret: "test-secret",
					},
				)
			},
			expectedError: nil,
			wantToken:     true,
		},
		{
			name: "empty login and password",
			input: models.UserCreateReq{
				Login:    "",
				Password: "",
			},
			mockSetup:     func(mus *MockUserStorage, mc *MockConfig) {},
			expectedError: errors.New("login and password are required"),
			wantToken:     false,
		},
		{
			name: "user not found",
			input: models.UserCreateReq{
				Login:    "nonexistent",
				Password: "password",
			},
			mockSetup: func(mus *MockUserStorage, mc *MockConfig) {
				pswdHash, _ := getHash("password")
				mus.On("Get", models.UserCreateDto{
					Login:        "nonexistent",
					PasswordHash: pswdHash,
				}).Return(
					(*models.UserCreateRes)(nil),
					nil,
				)
			},
			expectedError: ErrWrongLoginOrPassword,
			wantToken:     false,
		},
		{
			name: "error getting user",
			input: models.UserCreateReq{
				Login:    "erroruser",
				Password: "password",
			},
			mockSetup: func(mus *MockUserStorage, mc *MockConfig) {
				pswdHash, _ := getHash("password")
				mus.On("Get", models.UserCreateDto{
					Login:        "erroruser",
					PasswordHash: pswdHash,
				}).Return(
					(*models.UserCreateRes)(nil),
					errors.New("database error"),
				)
			},
			expectedError: errors.New("database error"),
			wantToken:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStorage := &MockUserStorage{}
			mockConfig := &MockConfig{}
			mockTokenStorage := &MockTokenStorage{}

			tt.mockSetup(mockUserStorage, mockConfig)

			// Очищаем контейнер перед каждым тестом
			testContainer = dig.New()
			app.AppContainer = testContainer

			err := testContainer.Provide(func() app.AppConfig {
				return mockConfig
			})
			assert.NoError(t, err)

			service := AuthNew(mockUserStorage, mockTokenStorage)

			token, err := service.Login(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotNil(t, token)
					assert.NotEmpty(t, token.Access)
					assert.NotEmpty(t, token.Refresh)
				} else {
					assert.Nil(t, token)
				}
			}

			mockUserStorage.AssertExpectations(t)
			mockConfig.AssertExpectations(t)
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	originalContainer := app.AppContainer
	defer func() { app.AppContainer = originalContainer }()

	testContainer := dig.New()
	app.AppContainer = testContainer

	// Создаем тестовый JWT токен для использования в тестах
	testUser := models.UserCreateRes{
		Id:    "test-user-id",
		Login: "testuser",
	}

	createTestToken := func(exp time.Time) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": testUser.Id,
			"exp": exp.Unix(),
			"iat": time.Now().Unix(),
		})
		signedToken, _ := token.SignedString([]byte("test-secret"))
		return signedToken
	}

	validRefreshToken := createTestToken(time.Now().Add(time.Hour))
	expiredRefreshToken := createTestToken(time.Now().Add(-time.Hour))

	tests := []struct {
		name           string
		refreshToken   string
		mockSetup      func(*MockUserStorage, *MockConfig, *MockTokenStorage)
		expectedError  error
		expectedTokens bool
	}{
		{
			name:         "successful token refresh",
			refreshToken: validRefreshToken,
			mockSetup: func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})
				mus.On("GetById", "test-user-id").Return(&testUser, nil)
				mts.On("SetTokens", mock.Anything, mock.Anything).Return(nil)
				mts.On("RemoveToken", mock.Anything, validRefreshToken, "").Return(nil)
			},
			expectedError:  nil,
			expectedTokens: true,
		},
		{
			name:           "empty refresh token",
			refreshToken:   "",
			mockSetup:      func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})

			},
			expectedError:  errors.New("invalid refresh token: token is malformed: token contains an invalid number of segments"),
			expectedTokens: false,
		},
		{
			name:         "expired refresh token",
			refreshToken: expiredRefreshToken,
			mockSetup: func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})
			},
			expectedError:  errors.New("invalid refresh token: token has invalid claims: token is expired"),
			expectedTokens: false,
		},
		{
			name:         "invalid token signature",
			refreshToken: "invalid.token.signature",
			mockSetup: func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})
			},
			expectedError:  errors.New("invalid refresh token: token is malformed: could not JSON decode header: invalid character '\\u008a' looking for beginning of value"),
			expectedTokens: false,
		},
		{
			name:         "user not found",
			refreshToken: validRefreshToken,
			mockSetup: func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})
				mus.On("GetById", "test-user-id").Return((*models.UserCreateRes)(nil), errors.New("user not found"))
			},
			expectedError:  errors.New("user not found: user not found"),
			expectedTokens: false,
		},
		{
			name:         "failed to store new tokens",
			refreshToken: validRefreshToken,
			mockSetup: func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})
				mus.On("GetById", "test-user-id").Return(&testUser, nil)
				mts.On("SetTokens", mock.Anything, mock.Anything).Return(errors.New("storage error"))
			},
			expectedError:  errors.New("failed to store tokens: storage error"),
			expectedTokens: false,
		},
		{
			name:         "success without token store",
			refreshToken: validRefreshToken,
			mockSetup: func(mus *MockUserStorage, mc *MockConfig, mts *MockTokenStorage) {
				mc.On("GetConfig").Return(&config.Config{Secret: "test-secret"})
				mus.On("GetById", "test-user-id").Return(&testUser, nil)
				mts.On("SetTokens", mock.Anything, mock.Anything).Return(nil)
				mts.On("RemoveToken", mock.Anything, validRefreshToken, "").Return(nil)

			},
			expectedError:  nil,
			expectedTokens: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStorage := &MockUserStorage{}
			mockConfig := &MockConfig{}
			mockTokenStorage := &MockTokenStorage{}

			tt.mockSetup(mockUserStorage, mockConfig, mockTokenStorage)

			// Очищаем контейнер перед каждым тестом
			testContainer = dig.New()
			app.AppContainer = testContainer

			err := testContainer.Provide(func() app.AppConfig {
				return mockConfig
			})
			assert.NoError(t, err)

			service := AuthNew(mockUserStorage, mockTokenStorage)
			service.tokenStore = mockTokenStorage

			tokens, err := service.RefreshToken(tt.refreshToken)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				if tt.expectedTokens {
					assert.NotNil(t, tokens)
					assert.NotEmpty(t, tokens.Access)
					assert.NotEmpty(t, tokens.Refresh)
				} else {
					assert.Nil(t, tokens)
				}
			}

			mockUserStorage.AssertExpectations(t)
			mockConfig.AssertExpectations(t)
			mockTokenStorage.AssertExpectations(t)
		})
	}
}
