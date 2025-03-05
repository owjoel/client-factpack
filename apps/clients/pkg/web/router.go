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
	"github.com/owjoel/client-factpack/apps/clients/config"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/owjoel/client-factpack/apps/clients/pkg/storage"
	"github.com/owjoel/client-factpack/apps/clients/pkg/web/handlers"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"}, // Allow frontend origin
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders: []string{"Content-Length"},
	}))
	

	pprof.Register(router)

	storage.Init()
	clientStorage := storage.GetInstance().Client
	clientService := service.NewClientService(clientStorage)
	handler := handlers.New(clientService)

	// Use RPC styling rather than REST
	v1API := router.Group("/api/v1/clients")
	v1API.GET("/health", handler.HealthCheck)
	v1API.GET("/retrieveProfile/:id", handler.GetClient)
	v1API.GET("/retrieveAllProfiles", handler.GetAllClients)
	v1API.POST("/createProfile", handler.CreateClient)

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
		log.Printf("started on port: %v\n", port)
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
