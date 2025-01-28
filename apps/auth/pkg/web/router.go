package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/auth"
	"github.com/owjoel/client-factpack/apps/auth/pkg/web/handlers"
)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	pprof.Register(router)

	handler := handlers.New()

	v1API := router.Group("/api/v1")
	v1API.GET("/health", handler.HealthCheck)
	v1API.POST("/createUser", handler.CreateUser)
	v1API.POST("/forgetPassword", handler.ForgetPassword)
	v1API.POST("/confirmForgetPassword", handler.ConfirmForgetPassword)
	v1API.POST("/user/login", handler.UserLogin)

	jwks, err := GetJWKS(config.AwsRegion, config.UserPoolId)
	if err != nil {
		log.Fatalf("Failed to retrieve Cognito JWKS\nError: %s", err)
	}

	v1API.GET("/protectedhealth", auth.AuthMiddleware("access", jwks), handler.HealthCheck)

	return &Router{router}
}

func (r *Router) Run() {
	port := config.GetPort(8080)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r.Engine,
	}

	go func() {
		err := srv.ListenAndServe()
		if err == http.ErrServerClosed {
			log.Fatalf("Server closed: %v\n", err)
		} else if err != nil {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

}

func Run() {
	NewRouter().Run()
}

func GetJWKS(awsRegion string, cognitoUserPoolId string) (*keyfunc.JWKS, error) {

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, cognitoUserPoolId)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}
	return jwks, nil
}
