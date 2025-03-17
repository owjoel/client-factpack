package service

import (
	"fmt"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
)

type APINotifier struct{}

// SendNotification sends a message to a user through WebSocket.
func (a *APINotifier) SendNotification(userID, message string) {
	fmt.Println("Sending notification to user:", userID)
	api.SendNotification(userID, message) // Send via WebSocket
}

type NotificationService struct {
	Notifier Notifier
}

// SendNotification calls the Notifier's SendNotification
func (s *NotificationService) SendNotification(userID, message string) {
	s.Notifier.SendNotification(userID, message)
}
