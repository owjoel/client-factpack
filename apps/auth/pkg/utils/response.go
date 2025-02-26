package utils

import (
    
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/owjoel/client-factpack/apps/auth/pkg/errors"
)

// ErrorResponse sends standardized error messages
func ErrorResponse(c *gin.Context, err errors.CustomError) {
    // Ensure that the error has a valid HTTP status
    if err.Status == 0 {
        err.Status = http.StatusInternalServerError // Default to 500 Internal Server Error
    }

    c.JSON(err.Status, gin.H{
        "error_code": err.Code,
        "message":    err.Message,
    })
    c.Abort()
}
