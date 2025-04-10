package api

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
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId query param is required"})
		return
	}

	notifications, err := Store.GetNotificationsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
