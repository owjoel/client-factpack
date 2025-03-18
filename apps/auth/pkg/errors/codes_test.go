package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	autherrors "github.com/owjoel/client-factpack/apps/auth/pkg/errors"
)

func TestGetError_KnownKeys(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		expectedError autherrors.CustomError
	}{
		{
			name:         "UserNotFound",
			key:          "UserNotFound",
			expectedError: autherrors.ErrUserNotFound,
		},
		{
			name:         "InvalidToken",
			key:          "InvalidToken",
			expectedError: autherrors.ErrInvalidToken,
		},
		{
			name:         "Unauthorized",
			key:          "Unauthorized",
			expectedError: autherrors.ErrUnauthorized,
		},
		{
			name:         "InvalidInput",
			key:          "InvalidInput",
			expectedError: autherrors.ErrInvalidInput,
		},
		{
			name:         "ClientNotFound",
			key:          "ClientNotFound",
			expectedError: autherrors.ErrClientNotFound,
		},
		{
			name:         "InternalError",
			key:          "InternalError",
			expectedError: autherrors.ErrServerError,
		},
		{
			name:         "WeakPassword",
			key:          "WeakPassword",
			expectedError: autherrors.ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := autherrors.GetError(tt.key)
			assert.Equal(t, tt.expectedError, result)
		})
	}
}

func TestGetError_UnknownKey_ReturnsServerError(t *testing.T) {
	key := "NonExistentKey"
	result := autherrors.GetError(key)
	assert.Equal(t, autherrors.ErrServerError, result)
}
