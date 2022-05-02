package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hahattan/assignment/application"
	"github.com/hahattan/assignment/interfaces/mocks"
)

func TestController_CreateToken(t *testing.T) {
	mockDB := &mocks.DBClient{}
	mockDB.On("AddToken", mock.Anything, mock.Anything).Return(nil)

	controller := NewController(mockDB)

	req, err := http.NewRequest(http.MethodPost, "/api/admin/token", http.NoBody)
	require.NoError(t, err)
	// Act
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.CreateToken)
	handler.ServeHTTP(recorder, req)

	var res string
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(res), 6)
	require.LessOrEqual(t, len(res), 12)
}

func TestController_DisableToken(t *testing.T) {
	validToken := "validToken"
	notFoundToken := "notFound"

	mockDB := &mocks.DBClient{}
	mockDB.On("TokenExists", validToken).Return(true, nil)
	mockDB.On("UpdateToken", validToken, mock.Anything).Return(nil)
	mockDB.On("TokenExists", notFoundToken).Return(false, nil)

	controller := NewController(mockDB)

	tests := []struct {
		name               string
		token              string
		expectedStatusCode int
	}{
		{"valid", validToken, http.StatusOK},
		{"invalid - token not found", notFoundToken, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/admin/token/disable", http.NoBody)
			req = mux.SetURLVars(req, map[string]string{"token": tt.token})
			require.NoError(t, err)
			// Act
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.DisableToken)
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, recorder.Result().StatusCode, tt.expectedStatusCode)
		})
	}
}

func TestController_GetAllTokens(t *testing.T) {
	mockData := map[string]string{
		"testToken1": time.Now().Format(application.TimeFormat),
		"testToken2": time.Time{}.Format(application.TimeFormat),
		"testToken3": time.Now().Add(time.Hour).Format(application.TimeFormat),
	}
	validData := map[string]string{
		"testToken1": "EXPIRED",
		"testToken2": "DISABLED",
		"testToken3": "VALID",
	}

	mockDB := &mocks.DBClient{}
	mockDB.On("AllTokens").Return(mockData, nil)

	controller := NewController(mockDB)

	req, err := http.NewRequest(http.MethodGet, "/api/admin/token", http.NoBody)
	require.NoError(t, err)
	// Act
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.GetAllTokens)
	handler.ServeHTTP(recorder, req)

	var res map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, res, validData)
}

func TestController_Login(t *testing.T) {
	validToken := "validToken"
	notFoundToken := "notFound"
	invalidTimeFormatToken := "wrongFormat"
	disabledToken := "disabled"
	expiredToken := "expired"

	mockDB := &mocks.DBClient{}
	mockDB.On("TokenExists", validToken).Return(true, nil)
	mockDB.On("TokenExists", notFoundToken).Return(false, nil)
	mockDB.On("TokenExists", invalidTimeFormatToken).Return(true, nil)
	mockDB.On("TokenExists", disabledToken).Return(true, nil)
	mockDB.On("TokenExists", expiredToken).Return(true, nil)
	mockDB.On("ExpiryByToken", validToken).Return(time.Now().Add(time.Minute).Format(application.TimeFormat), nil)
	mockDB.On("ExpiryByToken", invalidTimeFormatToken).Return("invalid", nil)
	mockDB.On("ExpiryByToken", disabledToken).Return(time.Time{}.Format(application.TimeFormat), nil)
	mockDB.On("ExpiryByToken", expiredToken).Return(time.Now().Format(application.TimeFormat), nil)

	controller := NewController(mockDB)

	tests := []struct {
		name               string
		token              interface{}
		expectedStatusCode int
	}{
		{"valid", validToken, http.StatusOK},
		{"valid - token disabled", disabledToken, http.StatusBadRequest},
		{"valid - token expired", expiredToken, http.StatusBadRequest},
		{"invalid - request token not string type", 123, http.StatusBadRequest},
		{"invalid - token not found", notFoundToken, http.StatusBadRequest},
		{"invalid - invalid time format", invalidTimeFormatToken, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.token)
			require.NoError(t, err)

			reader := strings.NewReader(string(jsonData))
			req, err := http.NewRequest(http.MethodPost, "/api/login", reader)
			require.NoError(t, err)
			// Act
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(controller.Login)
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, recorder.Result().StatusCode, tt.expectedStatusCode)
		})
	}
}