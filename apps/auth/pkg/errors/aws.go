package errors

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func isErrOfType[T error](err error) bool {
	var target T
	return errors.As(err, &target)
}

func CognitoErrorHandler(err error) (int, string) {
	switch {
	case isErrOfType[*types.UserNotFoundException](err):
		return http.StatusNotFound, "User not found"
	case isErrOfType[*types.PasswordResetRequiredException](err):
		return http.StatusForbidden, "Password reset required"
	case isErrOfType[*types.UserNotConfirmedException](err):
		return http.StatusForbidden, "User not confirmed"
	case isErrOfType[*types.InvalidParameterException](err):
		return http.StatusBadRequest, "Invalid parameters"
	case isErrOfType[*types.NotAuthorizedException](err):
		return http.StatusUnauthorized, "Incorrect username or password"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
