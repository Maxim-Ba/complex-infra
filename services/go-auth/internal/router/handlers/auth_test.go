package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-auth/internal/app"
	"go-auth/internal/models"
	"go-auth/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/dig"
)

type MockAuthService struct {
	mock.Mock
}

// Login implements app.AppAuthService.
func (m *MockAuthService) Login(user models.UserCreateReq) (*models.TokenDto, error) {
	args := m.Called(user)
	return args.Get(0).(*models.TokenDto), args.Error(1)
}

func (m *MockAuthService) Create(user models.UserCreateReq) (*models.TokenDto, error) {
	args := m.Called(user)
	return args.Get(0).(*models.TokenDto), args.Error(1)
}
func (m *MockAuthService) generateJWT(user models.UserCreateRes) (string, error) {
	args := m.Called(user)
	return "", args.Error(1)
}

type MockTokenStorage struct {
	mock.Mock
}

// RemoveToken implements app.AppTokenStorage.
func (m *MockTokenStorage) RemoveToken(ctx context.Context, refresh string, access string) error {
	args := m.Called(ctx, refresh, access)
	return args.Error(0)
}

func (m *MockTokenStorage) SetTokens(ctx context.Context, tokens *models.TokenDto) error {
	args := m.Called(ctx, tokens)
	return args.Error(0)
}

// Helper function to find cookie by name
func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

type RegisterTestCase struct {
	name            string
	requestBody     interface{}
	setupMocks      func(*MockAuthService, *MockTokenStorage)
	expectedStatus  int
	expectedBody    string
	checkCookies    bool
	expectedCookies map[string]string
}

func TestRegister(t *testing.T) {
	tests := []RegisterTestCase{
		{
			name: "Success",
			requestBody: models.UserCreateReq{
				Login:    "testuser",
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				tokens := &models.TokenDto{
					Access:  "access_token",
					Refresh: "refresh_token",
				}
				a.On("Create", mock.AnythingOfType("models.UserCreateReq")).Return(tokens, nil)
				ts.On("SetTokens", mock.Anything, tokens).Return(nil)
			},

			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
			checkCookies:   true,
			expectedCookies: map[string]string{
				"access_token":  "access_token",
				"refresh_token": "refresh_token",
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: "invalid json",
			setupMocks:  func(*MockAuthService, *MockTokenStorage) {},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
			checkCookies:   false,
		},
		{
			name: "MissingLoginOrPassword",
			requestBody: models.UserCreateReq{
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				a.On("Create", mock.AnythingOfType("models.UserCreateReq")).
					Return(&models.TokenDto{}, errors.New("login and password are required"))
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   "login and password are required",
			checkCookies:   false,
		},
		{
			name: "TokenStorageError",
			requestBody: models.UserCreateReq{
				Login:    "testuser",
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				tokens := &models.TokenDto{
					Access:  "access_token",
					Refresh: "refresh_token",
				}
				a.On("Create", mock.AnythingOfType("models.UserCreateReq")).Return(tokens, nil)
				ts.On("SetTokens", mock.Anything, tokens).Return(errors.New("storage error"))
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "token store error",
			checkCookies:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare mocks
			mockAuth := new(MockAuthService)
			mockTokenStorage := new(MockTokenStorage)
			tt.setupMocks(mockAuth, mockTokenStorage)

			// Prepare container with our mocks
			AppContainer := dig.New()
			AppContainer.Provide(func() app.AppAuthService { return mockAuth })
			AppContainer.Provide(func() app.AppTokenStorage { return mockTokenStorage })
			app.AppContainer = AppContainer

			// Create request
			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Call the handler
			Register(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			} else {
				assert.Empty(t, w.Body.String())
			}

			// Check cookies if needed
			if tt.checkCookies {
				cookies := w.Result().Cookies()
				for name, value := range tt.expectedCookies {
					cookie := findCookie(cookies, name)
					assert.NotNil(t, cookie, "cookie %s not found", name)
					assert.Equal(t, value, cookie.Value, "cookie %s value mismatch", name)
				}
			}

			// Verify all expectations were met
			mockAuth.AssertExpectations(t)
			mockTokenStorage.AssertExpectations(t)
		})
	}
}


type LogoutTestCase struct {
	name               string
	setupRequest       func(*http.Request)
	setupMocks         func(*MockTokenStorage)
	expectedStatus     int
	expectedBody       string
	checkCookiesCleared bool
}

func TestLogout(t *testing.T) {
	tests := []LogoutTestCase{
		{
			name: "Success",
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "access_token", Value: "valid_access"})
				r.AddCookie(&http.Cookie{Name: "refresh_token", Value: "valid_refresh"})
			},
			setupMocks: func(ts *MockTokenStorage) {
				ts.On("RemoveToken", mock.Anything, "valid_refresh", "valid_access").Return(nil)
			},
			expectedStatus:     http.StatusNoContent,
			expectedBody:       "",
			checkCookiesCleared: true,
		},
		{
			name: "MissingAccessToken",
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "refresh_token", Value: "valid_refresh"})
			},
			setupMocks:         func(*MockTokenStorage) {},
			expectedStatus:     http.StatusInternalServerError,
			expectedBody:       "Failed to read cookie access",
			checkCookiesCleared: false,
		},
		{
			name: "MissingRefreshToken",
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "access_token", Value: "valid_access"})
			},
			setupMocks:         func(*MockTokenStorage) {},
			expectedStatus:     http.StatusInternalServerError,
			expectedBody:       "Failed to read cookie refresh",
			checkCookiesCleared: false,
		},
		{
			name: "TokenStorageError",
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "access_token", Value: "valid_access"})
				r.AddCookie(&http.Cookie{Name: "refresh_token", Value: "valid_refresh"})
			},
			setupMocks: func(ts *MockTokenStorage) {
				ts.On("RemoveToken", mock.Anything, "valid_refresh", "valid_access").
					Return(errors.New("storage error"))
			},
			expectedStatus:     http.StatusInternalServerError,
			expectedBody:       "Failed to remove token",
			checkCookiesCleared: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare mocks
			mockTokenStorage := new(MockTokenStorage)
			tt.setupMocks(mockTokenStorage)

			// Prepare container with our mock
			AppContainer := dig.New()
			AppContainer.Provide(func() app.AppTokenStorage { return mockTokenStorage })
			app.AppContainer = AppContainer

			// Create request
			req := httptest.NewRequest("POST", "/logout", nil)
			tt.setupRequest(req)
			w := httptest.NewRecorder()

			// Call the handler
			Logout(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			} else {
				assert.Empty(t, w.Body.String())
			}

			// Check cookies are cleared if needed
			if tt.checkCookiesCleared {
				cookies := w.Result().Cookies()
				accessCookie := findCookie(cookies, "access_token")
				refreshCookie := findCookie(cookies, "refresh_token")
				
				assert.NotNil(t, accessCookie, "access_token cookie not found")
				assert.NotNil(t, refreshCookie, "refresh_token cookie not found")
				assert.Equal(t, "", accessCookie.Value, "access_token should be cleared")
				assert.Equal(t, "", refreshCookie.Value, "refresh_token should be cleared")
				assert.Equal(t, -1, accessCookie.MaxAge, "access_token should be expired")
				assert.Equal(t, -1, refreshCookie.MaxAge, "refresh_token should be expired")
			}

			// Verify all expectations were met
			mockTokenStorage.AssertExpectations(t)
		})
	}
}

type LoginTestCase struct {
	name            string
	requestBody     interface{}
	setupMocks      func(*MockAuthService, *MockTokenStorage)
	expectedStatus  int
	expectedBody    string
	checkCookies    bool
	expectedCookies map[string]string
}

func TestLogin(t *testing.T) {
	tests := []LoginTestCase{
		{
			name: "Success",
			requestBody: models.UserCreateReq{
				Login:    "testuser",
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				tokens := &models.TokenDto{
					Access:  "access_token",
					Refresh: "refresh_token",
				}
				a.On("Login", mock.AnythingOfType("models.UserCreateReq")).Return(tokens, nil)
				ts.On("SetTokens", mock.Anything, tokens).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
			checkCookies:   true,
			expectedCookies: map[string]string{
				"access_token":  "access_token",
				"refresh_token": "refresh_token",
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: "invalid json",
			setupMocks:  func(*MockAuthService, *MockTokenStorage) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
			checkCookies:   false,
		},
		{
			name: "MissingLoginOrPassword",
			requestBody: models.UserCreateReq{
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				a.On("Login", mock.AnythingOfType("models.UserCreateReq")).
					Return(&models.TokenDto{}, services.ErrLoginAndPasswordAreRequired)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Wrong login or password",
			checkCookies:   false,
		},
		{
			name: "WrongLoginOrPassword",
			requestBody: models.UserCreateReq{
				Login:    "wronguser",
				Password: "wrongpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				a.On("Login", mock.AnythingOfType("models.UserCreateReq")).
					Return(&models.TokenDto{}, services.ErrWrongLoginOrPassword)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Wrong login or password",
			checkCookies:   false,
		},
		{
			name: "TokenStorageError",
			requestBody: models.UserCreateReq{
				Login:    "testuser",
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				tokens := &models.TokenDto{
					Access:  "access_token",
					Refresh: "refresh_token",
				}
				a.On("Login", mock.AnythingOfType("models.UserCreateReq")).Return(tokens, nil)
				ts.On("SetTokens", mock.Anything, tokens).Return(errors.New("storage error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "token store error",
			checkCookies:   false,
		},
		{
			name: "AuthServiceError",
			requestBody: models.UserCreateReq{
				Login:    "testuser",
				Password: "testpass",
			},
			setupMocks: func(a *MockAuthService, ts *MockTokenStorage) {
				a.On("Login", mock.AnythingOfType("models.UserCreateReq")).
					Return(&models.TokenDto{}, errors.New("some auth error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad request",
			checkCookies:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare mocks
			mockAuth := new(MockAuthService)
			mockTokenStorage := new(MockTokenStorage)
			tt.setupMocks(mockAuth, mockTokenStorage)

			// Prepare container with our mocks
			AppContainer := dig.New()
			AppContainer.Provide(func() app.AppAuthService { return mockAuth })
			AppContainer.Provide(func() app.AppTokenStorage { return mockTokenStorage })
			app.AppContainer = AppContainer

			// Create request
			var body []byte
			switch v := tt.requestBody.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Call the handler
			Login(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			} else {
				assert.Empty(t, w.Body.String())
			}

			// Check cookies if needed
			if tt.checkCookies {
				cookies := w.Result().Cookies()
				for name, value := range tt.expectedCookies {
					cookie := findCookie(cookies, name)
					assert.NotNil(t, cookie, "cookie %s not found", name)
					assert.Equal(t, value, cookie.Value, "cookie %s value mismatch", name)
				}
			}

			// Verify all expectations were met
			mockAuth.AssertExpectations(t)
			mockTokenStorage.AssertExpectations(t)
		})
	}
}
