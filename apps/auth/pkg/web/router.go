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

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/auth/config"
	"github.com/owjoel/client-factpack/apps/auth/pkg/web/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router represents the router for the web service.
type Router struct {
	*gin.Engine
}

// NewRouter creates a new router.
func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// enable CORS
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		}),
	)

	pprof.Register(router)

	handler := handlers.New()

	// Use RPC styling rather than REST
	v1API := router.Group("/api/v1")
	v1API.GET("/health", handler.HealthCheck)
	{
		auth := v1API.Group("/auth")
		auth.POST("/createUser", handler.CreateUser)
		auth.POST("/forgetPassword", handler.ForgetPassword)
		auth.POST("/confirmForgetPassword", handler.ConfirmForgetPassword)

		auth.POST("/login", handler.UserLogin)
		auth.POST("/changePassword", handler.UserInitialChangePassword)
		auth.GET("/setupMFA", handler.UserSetupMFA)
		auth.POST("/verifyMFA", handler.UserVerifyMFA)
		auth.POST("/loginMFA", handler.UserLoginMFA)
		auth.POST("/logout", handler.UserLogout)

		auth.Use(handler.Authenticate)
		auth.GET("/checkUser", handler.HealthCheck)
		auth.GET("/username", handler.GetUsername)
		auth.GET("/role", handler.GetUserRole)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &Router{router}
}

// Run starts the web service.
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

// Run starts the web service.
func Run() {
	NewRouter().Run()
}
