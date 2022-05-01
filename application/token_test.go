package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hahattan/assignment/interfaces/mocks"
)

func TestCreateToken(t *testing.T) {
	mockDB := &mocks.DBClient{}
	mockDB.On("AddToken", mock.Anything, mock.Anything).Return(nil)

	res, err := CreateToken(mockDB)
	require.NoError(t, err)
	require.LessOrEqual(t, len(res), 12)
	require.GreaterOrEqual(t, len(res), 6)
}

func TestDisableToken(t *testing.T) {
	testToken := "test123"
	mockDB := &mocks.DBClient{}
	mockDB.On("UpdateToken", testToken, time.Now().Format(timeFormat)).Return(nil)

	err := DisableToken(testToken, mockDB)
	require.NoError(t, err)
}

func TestGetAllTokens(t *testing.T) {
	mockData := map[string]string{
		"testToken1": time.Now().Format(timeFormat),
		"testToken2": time.Time{}.Format(timeFormat),
		"testToken3": time.Now().Add(time.Hour).Format(timeFormat),
	}
	validData := map[string]string{
		"testToken1": "EXPIRED",
		"testToken2": "DISABLED",
		"testToken3": "VALID",
	}
	mockDB := &mocks.DBClient{}
	mockDB.On("AllTokens").Return(mockData, nil)

	res, err := GetAllTokens(mockDB)
	require.NoError(t, err)
	require.Equal(t, res, validData)
}

func TestValidateToken(t *testing.T) {
	testToken := "testToken"
	invalidTimeFormatToken := "invalidTimeFormatToken"
	disabledToken := "disabledToken"
	expiredToken := "expiredToken"

	mockDB := &mocks.DBClient{}
	mockDB.On("ExpiryByToken", testToken).Return(time.Now().Add(time.Minute).Format(timeFormat), nil)
	mockDB.On("ExpiryByToken", invalidTimeFormatToken).Return("invalid", nil)
	mockDB.On("ExpiryByToken", disabledToken).Return(time.Time{}.Format(timeFormat), nil)
	mockDB.On("ExpiryByToken", expiredToken).Return(time.Now().Format(timeFormat), nil)

	tests := []struct {
		name        string
		token       string
		expectedErr bool
	}{
		{"valid", testToken, false},
		{"valid - token disabled", disabledToken, true},
		{"valid - token expired", expiredToken, true},
		{"invalid - invalid time format", invalidTimeFormatToken, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token, mockDB)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
