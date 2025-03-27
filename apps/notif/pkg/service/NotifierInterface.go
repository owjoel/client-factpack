package service

// Notifier defines the behavior for sending notifications.
type Notifier interface {
	SendNotification(userID, message string)
}
