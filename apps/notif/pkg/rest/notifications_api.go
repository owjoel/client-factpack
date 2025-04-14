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

// GET /api/v1/notifications?username=xxx
func GetUserNotifications(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username query param is required"})
		return
	}

	status := c.Query("status")
	page := parseIntWithDefault(c.Query("page"), 1)
	pageSize := parseIntWithDefault(c.Query("pageSize"), 10)

	notifications, err := Store.GetNotificationsByUser(username, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"page":          page,
		"pageSize":      pageSize,
	})
}


func GetClientNotifications(c *gin.Context) {
	name := c.Query("name")
	priority := c.Query("priority")
	page := parseIntWithDefault(c.Query("page"), 1)
	pageSize := parseIntWithDefault(c.Query("pageSize"), 10)

	notifications, err := Store.GetClientNotifications(name, priority, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve client notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"page":          page,
		"pageSize":      pageSize,
	})
}

