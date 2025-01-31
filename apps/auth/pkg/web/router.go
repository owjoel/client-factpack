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

type Router struct {
	*gin.Engine
}

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
		auth.POST("/login", handler.UserLogin) // maybe wanna group admin endpoints together
	}

	// MFA
	v1API.POST("/auth/associateToken", handler.AssociateToken)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
