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

// GetUserNotifications godoc
// @Summary      Get job notifications
// @Description  Returns job notifications for a given username, optionally filtered by status
// @Tags         notifications
// @Param        username  query  string  true  "Username"
// @Param        status    query  string  false "Job status filter (completed, pending, failed, processing)"
// @Param        page      query  int     false "Page number"
// @Param        pageSize  query  int     false "Number of items per page"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications [get]
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


// GetClientNotifications godoc
// @Summary      Get client notifications
// @Description  Returns paginated client notifications, optionally filtered by name and priority
// @Tags         notifications
// @Param        name      query  string  false "Client name filter"
// @Param        priority  query  string  false "Priority filter (low, medium, high)"
// @Param        page      query  int     false "Page number"
// @Param        pageSize  query  int     false "Number of items per page"
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /notifications/client [get]
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

