package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

// Authenticate is a middleware that checks if the user is authenticated by validating the "accessToken" cookie.
// If the cookie is missing or invalid, it returns a 403 Forbidden response.
// Otherwise, it sets the token in the context for downstream handlers to use.
func (h *UserHandler) Authenticate(c *gin.Context) {
	token, err := c.Cookie("accessToken")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Not signed in"})
		c.Abort()
		return
	}

	c.Set("accessToken", token)
	c.Next()
}

// VerifyMFA verifies the user's multi-factor authentication (MFA) setup using their access token.
// It retrieves the "accessToken" from the context, and if it doesn't exist, returns a 403 Forbidden response.
// If the token exists, it prepares the input for associating an MFA software token with AWS Cognito.
func (h *UserHandler) VerifyMFA(c *gin.Context) {
	token, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"message": "Could not verify identity"})
	}
	tokenString := fmt.Sprintf("%v", token)
	input := &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(tokenString),
	}
	fmt.Println(input)
}
