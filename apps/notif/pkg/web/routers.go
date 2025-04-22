package web

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/owjoel/client-factpack/apps/notif/docs"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
	"github.com/owjoel/client-factpack/apps/notif/pkg/rest"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Initializes Gin router
func InitRouter() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// enable CORS
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173", "http://localhost:4173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		}),
	)

	// Health Check
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		api.HandleWebSocketConnections(c.Writer, c.Request)
	})

	// REST: Get notifications
	router.GET("/api/v1/notif/jobs", rest.GetUserNotifications)
	router.GET("/api/v1/notif/clients", rest.GetClientNotifications)
	port := ":8082"
	fmt.Println("Starting server on port", port)
	router.Run(port)
}
