package main

import (
	"net/http"

	"github.com/owjoel/client-factpack/apps/notif/pkg/web"
	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
    "github.com/owjoel/client-factpack/apps/notif/config"
    "github.com/owjoel/client-factpack/apps/notif/pkg/utils"
)

// Swagger
// @title		client-factpack/notifications
// @version	1.0
// @description	Notification service for handling real-time WebSocket notifications
// @host		localhost:8081
// @BasePath	/api/v1

func main() {
	utils.InitLogger() // Initialize logger
	utils.Logger.Info("Starting WebSocket Notification Service...")
    config.Load() 
	storage.InitMessageQueue() // Initialize message queue subscription
	web.InitRouter() // Set up WebSocket routing

	port := ":8081"
	utils.Logger.Infof("Listening on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		utils.Logger.Fatal("Server error:", err)
	}
}