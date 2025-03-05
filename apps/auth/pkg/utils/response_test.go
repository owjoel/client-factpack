package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Test case for a valid error response
func TestErrorResponse_Success(t *testing.T) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a new gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Define a sample error
	testErr := errors.CustomError{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_INPUT",
		Message: "Invalid request data",
	}

	// Call the function under test
	ErrorResponse(c, testErr)

	// Check response status
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Expected JSON response
	expectedResponse := `{"error_code":"INVALID_INPUT","message":"Invalid request data"}`

	// Check response body
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestErrorResponse_Fail(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    // Define an invalid error (empty CustomError struct)
    testErr := errors.CustomError{}

    // Call the function under test
    ErrorResponse(c, testErr)

    // Ensure it doesn't return 200 OK
    assert.NotEqual(t, http.StatusOK, w.Code, "Expected non-200 status but got 200")

    // Log response for debugging
    t.Logf("Response Status: %d, Body: %s", w.Code, w.Body.String())
}

