package service

import (
	"fmt"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
)

// Sends a notification
func SendNotification(userID, message string) {
	fmt.Println("Sending notification to user:", userID)
	api.SendNotification(userID, message) // Send via WebSocket
}
