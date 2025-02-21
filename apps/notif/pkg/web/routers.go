package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
)

// Initializes Gin router
func InitRouter() {
	router := gin.Default()

	// Health Check
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		api.HandleWebSocketConnections(c.Writer, c.Request)
	})

	port := ":8081"
	fmt.Println("Starting server on port", port)
	router.Run(port)
}
