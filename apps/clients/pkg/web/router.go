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
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/owjoel/client-factpack/apps/clients/pkg/web/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// enable CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))

	pprof.Register(router)

	mongoDb := repository.InitMongo()

	logRepository := repository.NewMongoLogRepository(mongoDb)
	logService := service.NewLogService(logRepository)
	logHandler := handlers.NewLogHandler(logService)

	jobRepository := repository.NewMongoJobRepository(mongoDb)
	jobService := service.NewJobService(jobRepository)
	jobHandler := handlers.NewJobHandler(jobService)

	clientRepository := repository.NewMongoClientRepository(mongoDb)
	clientService := service.NewClientService(clientRepository, jobService, logService)
	clientHandler := handlers.NewClientHandler(clientService)

	// Use RPC styling rather than REST

	// startregion Clients
	v1API := router.Group("/api/v1/clients")
	v1API.GET("/health", clientHandler.HealthCheck)
	v1API.GET("/:id", clientHandler.GetClient)
	v1API.GET("/", clientHandler.GetAllClients)
	v1API.PUT("/:id", clientHandler.UpdateClient)
	v1API.POST("/scrape", clientHandler.CreateClientByName)
	v1API.POST("/:id/scrape", clientHandler.RescrapeClient)
	v1API.POST("/match", clientHandler.MatchClient)
	// endregion Clients

	// startregion Jobs
	v1Jobs := router.Group("/api/v1/jobs")
	v1Jobs.GET("/:id", jobHandler.GetJob)
	v1Jobs.GET("/", jobHandler.GetAllJobs)
	// endregion Jobs

	// startregion Logs
	v1Logs := router.Group("/api/v1/logs")
	v1Logs.GET("/", logHandler.GetLogs)
	v1Logs.GET("/:id", logHandler.GetLog)
	// endregion Logs

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
