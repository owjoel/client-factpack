package rest

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
)

var Store *storage.NotificationStorage

// Inject DB store
func InitNotificationAPI(store *storage.NotificationStorage) {
	Store = store
}

// GET /api/v1/notifications?userId=xxx
func GetUserNotifications(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username query param is required"})
		return
	}

	notifications, err := Store.GetNotificationsByUser(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
