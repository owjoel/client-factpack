package handlers

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

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

func (h *UserHandler) AssociateToken(c *gin.Context) {
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