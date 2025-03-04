package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/errors"
	"github.com/owjoel/client-factpack/apps/auth/pkg/utils"
)

// Authenticate is a middleware that checks if the user is authenticated by validating the "accessToken" cookie.
func (h *UserHandler) Authenticate(c *gin.Context) {

	requiredTokenUse := "access" // default check for access token
	awsDefaultRegion := config.AwsRegion
	cognitoUserPoolId := config.UserPoolID
	cognitoAppClientId := config.ClientID

	jwks, err := GetJWKS(awsDefaultRegion, cognitoUserPoolId)
	if err != nil {
		log.Printf("Failed to retrieve Cognito JWKS: %s", err)
		utils.ErrorResponse(c, errors.ErrServerError)
		return
	}

	tokenString, err := c.Cookie("access_token")
	if err != nil || tokenString == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// * Verify the signature of the JWT
	// * Verify that the algorithm used is RS256
	// * Verify that the 'exp' claim exists in the token
	// * Verification of audience 'aud' is taken care later when we examine if the
	//   token is 'id' or 'access'
	// * The issuer (iss) claim should match your user pool. For example, a user
	//   pool created in the us-east-1 region
	//   will have the following iss value: https://cognito-idp.us-east-1.amazonaws.com/<userpoolID>.
	token, err := jwt.Parse(tokenString,
		jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", awsDefaultRegion, cognitoUserPoolId)))
	if err != nil || !token.Valid {
		utils.ErrorResponse(c, errors.ErrInvalidToken)
		return
	}

	// Parse JWT claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.ErrorResponse(c, errors.ErrInvalidToken)
		return
	}

	// Validate token expiration
	expClaim, err := claims.GetExpirationTime()
	if err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidToken)
		return
	}
	if expClaim.Unix() < time.Now().Unix() {
		utils.ErrorResponse(c, errors.ErrInvalidToken)
		return
	}

	// Validate token use
	tokenUseClaim, ok := claims["token_use"].(string)
	if !ok || tokenUseClaim != requiredTokenUse {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Validate subject claim (user identifier)
	subClaim, err := claims.GetSubject()
	if err != nil {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}
	c.Set("username", subClaim)

	// Extract the username from claims
	usernameClaim, ok := claims["username"].(string)
	if !ok {
		usernameClaim = subClaim // fallback to subject if username is missing
	}

	c.Set("username", usernameClaim) // Store non-hashed username


	// Validate client ID
	var appClientIdClaim string
	if tokenUseClaim == "id" {
		audienceClaims, err := claims.GetAudience()
		if err != nil || len(audienceClaims) == 0 {
			utils.ErrorResponse(c, errors.ErrUnauthorized)
			return
		}
		appClientIdClaim = audienceClaims[0]
	} else if tokenUseClaim == "access" {
		clientIdClaim, ok := claims["client_id"].(string)
		if !ok {
			utils.ErrorResponse(c, errors.ErrUnauthorized)
			return
		}
		appClientIdClaim = clientIdClaim
	} else {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	if appClientIdClaim != cognitoAppClientId {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Retrieve Cognito user groups
	userGroupsAttribute, ok := claims["cognito:groups"]
	userGroupsClaims := make([]string, 0)
	if ok {
		switch x := userGroupsAttribute.(type) {
		case []interface{}:
			for _, e := range x {
				userGroupsClaims = append(userGroupsClaims, e.(string))
			}
		default:
			utils.ErrorResponse(c, errors.ErrUnauthorized)
			return
		}
	}
	c.Set("groups", userGroupsClaims)

	c.Next()
	c.Set("accessToken", token)
	c.Next()
}

// VerifyMFA verifies the user's multi-factor authentication (MFA) setup using their access token.
func (h *UserHandler) VerifyMFA(c *gin.Context) {
	token, exists := c.Get("accessToken")
	if !exists {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	tokenString := fmt.Sprintf("%v", token)
	input := &cognitoidentityprovider.AssociateSoftwareTokenInput{
		AccessToken: aws.String(tokenString),
	}
	log.Println(input)
}

// GetJWKS retrieves Cognito JSON Web Key Set (JWKS) for verifying JWTs.
func GetJWKS(awsRegion string, cognitoUserPoolId string) (*keyfunc.JWKS, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, cognitoUserPoolId)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}
	return jwks, nil
}

func (h *UserHandler) GetUsername(c *gin.Context) {
    username, exists := c.Get("username")
    if !exists {
        utils.ErrorResponse(c, errors.ErrUnauthorized)
        return
    }

    c.JSON(200, gin.H{
        "username": username, // Return actual username
    })
}

// GetUserRole extracts the user's role from the JWT token
func (h *UserHandler) GetUserRole(c *gin.Context) {
	tokenString, err := c.Cookie("access_token")

	if err != nil || tokenString == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Call the service layer to get the user's role
	role, err := h.service.GetUserRoleFromToken(tokenString)
	if err != nil {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Return user role in JSON response
	c.JSON(200, gin.H{"user_role": role})
}