package errors_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/assert"

	autherrors "github.com/owjoel/client-factpack/apps/auth/pkg/errors"
)

func TestCognitoErrorHandler(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		expectedError autherrors.CustomError
	}{
		{
			name:         "UserNotFoundException",
			inputError:   &types.UserNotFoundException{},
			expectedError: autherrors.ErrUserNotFound,
		},
		{
			name:         "PasswordResetRequiredException",
			inputError:   &types.PasswordResetRequiredException{},
			expectedError: autherrors.CustomError{
				Code:    "AUTH_PASSWORD_RESET_REQUIRED",
				Message: "Password reset required",
				Status:  403,
			},
		},
		{
			name:         "UserNotConfirmedException",
			inputError:   &types.UserNotConfirmedException{},
			expectedError: autherrors.CustomError{
				Code:    "AUTH_USER_NOT_CONFIRMED",
				Message: "User not confirmed",
				Status:  403,
			},
		},
		{
			name:         "InvalidParameterException",
			inputError:   &types.InvalidParameterException{},
			expectedError: autherrors.ErrInvalidInput,
		},
		{
			name:         "NotAuthorizedException",
			inputError:   &types.NotAuthorizedException{},
			expectedError: autherrors.ErrUnauthorized,
		},
		{
			name:         "InvalidPasswordException",
			inputError:   &types.InvalidPasswordException{},
			expectedError: autherrors.ErrWeakPassword,
		},
		{
			name:         "UnknownError_DefaultServerError",
			inputError:   errors.New("some unknown error"),
			expectedError: autherrors.ErrServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := autherrors.CognitoErrorHandler(tt.inputError)
			assert.Equal(t, tt.expectedError, result)
		})
	}
}
