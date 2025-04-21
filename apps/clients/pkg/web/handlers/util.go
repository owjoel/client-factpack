package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/clients/config"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
)

type Response struct {
	ApiVersion string      `json:"version"`
	Timestamp  time.Time   `json:"timestamp"`
	Status     int         `json:"status"`
	Data       interface{} `json:"data"`
}

func resp(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, Response{
		ApiVersion: config.GetVersion(),
		Timestamp:  time.Now(),
		Status:     code,
		Data:       obj,
	})
}


func ErrorHandler(c *gin.Context, err error, message string) {
	switch {
	case errors.Is(err, errorx.ErrBadRequest):
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Bad request: " + message})
	case errors.Is(err, errorx.ErrInvalidInput):
		resp(c, http.StatusBadRequest, model.ErrorResponse{Message: "Invalid input: " + message})
	case errors.Is(err, errorx.ErrValidationFailed):
		resp(c, http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Validation failed: " + message})
	case errors.Is(err, errorx.ErrUnauthorized):
		resp(c, http.StatusUnauthorized, model.ErrorResponse{Message: "Unauthorized: " + message})
	case errors.Is(err, errorx.ErrForbidden):
		resp(c, http.StatusForbidden, model.ErrorResponse{Message: "Forbidden: " + message})
	case errors.Is(err, errorx.ErrNotFound):
		resp(c, http.StatusNotFound, model.ErrorResponse{Message: "Not found: " + message})
	case errors.Is(err, errorx.ErrConflict):
		resp(c, http.StatusConflict, model.ErrorResponse{Message: "Conflict: " + message})
	case errors.Is(err, errorx.ErrInternal):
		resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Internal server error: " + message})
	case errors.Is(err, errorx.ErrDependencyFailed):
		resp(c, http.StatusBadGateway, model.ErrorResponse{Message: "Upstream service failed: " + message})
	case errors.Is(err, errorx.ErrTimeout):
		resp(c, http.StatusGatewayTimeout, model.ErrorResponse{Message: "Operation timed out: " + message})
	default:
		resp(c, http.StatusInternalServerError, model.ErrorResponse{Message: "Unexpected error: " + message})
	}
}