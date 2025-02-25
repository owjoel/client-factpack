package utils

import (
    "github.com/gin-gonic/gin"
    "github.com/owjoel/client-factpack/apps/auth/pkg/errors"
)

// ErrorResponse sends standardized error messages
func ErrorResponse(c *gin.Context, err errors.CustomError) {
    c.JSON(err.Status, gin.H{
        "error_code": err.Code,
        "message":    err.Message,
    })
    c.Abort()
}

