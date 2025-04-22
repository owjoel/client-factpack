package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owjoel/client-factpack/apps/clients/config"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
)

// Authenticate is a middleware that checks if the user is authenticated by validating the "accessToken" cookie.
func Authenticate(getJWKS func(string, string) (*keyfunc.JWKS, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		requiredTokenUse := "access" // default check for access token
		awsDefaultRegion := config.AwsRegion
		cognitoUserPoolId := config.UserPoolID
		cognitoAppClientId := config.ClientID

		jwks, err := getJWKS(awsDefaultRegion, cognitoUserPoolId)
		if err != nil {
			log.Printf("Failed to retrieve Cognito JWKS: %s", err)
			resp(c, http.StatusInternalServerError, errorx.ErrInternal.Error())
			c.Abort()
			return
		}

		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			log.Printf("No access token found in cookie: %s", err)
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
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
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		// Parse JWT claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		// Validate token expiration
		expClaim, err := claims.GetExpirationTime()
		if err != nil {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}
		if expClaim.Unix() < time.Now().Unix() {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		// Validate token use
		tokenUseClaim, ok := claims["token_use"].(string)
		if !ok || tokenUseClaim != requiredTokenUse {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		// Validate subject claim (user identifier)
		subClaim, err := claims.GetSubject()
		if err != nil {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}
		c.Set("username", subClaim)

		// Extract the username from claims
		usernameClaim, ok := claims["username"].(string)
		if !ok {
			usernameClaim = subClaim // fallback to subject if username is missing
		}

		c.Set("username", usernameClaim) // Store non-hashed username

		// for context.Context to use in service layer
		ctx := context.WithValue(c.Request.Context(), "username", usernameClaim)
		c.Request = c.Request.WithContext(ctx)

		// Validate client ID
		var appClientIdClaim string
		if tokenUseClaim == "id" {
			audienceClaims, err := claims.GetAudience()
			if err != nil || len(audienceClaims) == 0 {
				resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
				c.Abort()
				return
			}
			appClientIdClaim = audienceClaims[0]
		} else if tokenUseClaim == "access" {
			clientIdClaim, ok := claims["client_id"].(string)
			if !ok {
				resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
				c.Abort()
				return
			}
			appClientIdClaim = clientIdClaim
		} else {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
			return
		}

		if appClientIdClaim != cognitoAppClientId {
			resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
			c.Abort()
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
				resp(c, http.StatusUnauthorized, errorx.ErrUnauthorized.Error())
				c.Abort()
				return
			}
		}
		c.Set("groups", userGroupsClaims)

		c.Next()
		c.Set("accessToken", token)
		c.Next()
	}
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
