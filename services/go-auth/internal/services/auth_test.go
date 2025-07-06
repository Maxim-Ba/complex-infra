package services

import (
	"errors"
	"testing"

	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/dig"
)

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) Save(user models.UserCreateReq) models.UserCreateRes {
	args := m.Called(user)
	return args.Get(0).(models.UserCreateRes)
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
					},
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
			expectedError: errors.New("database error"),
			wantToken:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStorage := &MockUserStorage{}
			mockConfig := &MockConfig{}

			tt.mockSetup(mockUserStorage, mockConfig)

			// Очищаем контейнер перед каждым тестом
			testContainer = dig.New()
			app.AppContainer = testContainer

			err := testContainer.Provide(func() app.AppConfig {
				return mockConfig
			})
			assert.NoError(t, err)

			service := AuthNew(mockUserStorage)

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
