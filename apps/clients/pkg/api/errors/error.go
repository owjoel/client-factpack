package errorx

import "errors"

var (
	// 400 errors
	ErrBadRequest       = errors.New("bad request")              // 400
	ErrInvalidInput     = errors.New("invalid input")            // 400
	ErrValidationFailed = errors.New("validation failed")        // 422
	ErrUnauthorized     = errors.New("unauthorized")             // 401
	ErrForbidden        = errors.New("forbidden")                // 403
	ErrNotFound         = errors.New("not found")                // 404
	ErrConflict         = errors.New("conflict")                 // 409

	// 500 errors
	ErrInternal         = errors.New("internal server error")    // 500
	ErrDependencyFailed = errors.New("upstream service failed")  // 502
	ErrTimeout          = errors.New("operation timed out")      // 504
)
