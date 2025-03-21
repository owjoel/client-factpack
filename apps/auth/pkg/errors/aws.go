package errors

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

// isErrOfType checks if an error matches a given AWS Cognito error type
func isErrOfType[T error](err error) bool {
	var target T
	return errors.As(err, &target)
}

// CognitoErrorHandler maps AWS Cognito errors to our CustomError structure
func CognitoErrorHandler(err error) CustomError {
	switch {
	case isErrOfType[*types.UserNotFoundException](err):
		return ErrUserNotFound
	case isErrOfType[*types.PasswordResetRequiredException](err):
		return CustomError{"AUTH_PASSWORD_RESET_REQUIRED", "Password reset required", 403}
	case isErrOfType[*types.UserNotConfirmedException](err):
		return CustomError{"AUTH_USER_NOT_CONFIRMED", "User not confirmed", 403}
	case isErrOfType[*types.InvalidParameterException](err):
		return ErrInvalidInput
	case isErrOfType[*types.NotAuthorizedException](err):
		return ErrUnauthorized
	case isErrOfType[*types.InvalidPasswordException](err):
		return ErrWeakPassword
	default:
		return ErrServerError
	}
}