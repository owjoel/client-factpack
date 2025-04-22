package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/web/handlers"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		err        error
		wantCode   int
		wantBody   string
	}{
		{"BadRequest", errorx.ErrBadRequest, http.StatusBadRequest, "Bad request: test message"},
		{"InvalidInput", errorx.ErrInvalidInput, http.StatusBadRequest, "Invalid input: test message"},
		{"ValidationFailed", errorx.ErrValidationFailed, http.StatusUnprocessableEntity, "Validation failed: test message"},
		{"Unauthorized", errorx.ErrUnauthorized, http.StatusUnauthorized, "Unauthorized: test message"},
		{"Forbidden", errorx.ErrForbidden, http.StatusForbidden, "Forbidden: test message"},
		{"NotFound", errorx.ErrNotFound, http.StatusNotFound, "Not found: test message"},
		{"Conflict", errorx.ErrConflict, http.StatusConflict, "Conflict: test message"},
		{"Internal", errorx.ErrInternal, http.StatusInternalServerError, "Internal server error: test message"},
		{"DependencyFailed", errorx.ErrDependencyFailed, http.StatusBadGateway, "Upstream service failed: test message"},
		{"Timeout", errorx.ErrTimeout, http.StatusGatewayTimeout, "Operation timed out: test message"},
		{"Default", errors.New("something else"), http.StatusInternalServerError, "Unexpected error: test message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handlers.ErrorHandler(c, tt.err, "test message")

			if w.Code != tt.wantCode {
				t.Errorf("Expected status code %d, got %d", tt.wantCode, w.Code)
			}

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			if err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			data := resp["data"].(map[string]interface{})
			if data["message"] != tt.wantBody {
				t.Errorf("Expected message '%s', got '%s'", tt.wantBody, data["message"])
			}
		})
	}
}
