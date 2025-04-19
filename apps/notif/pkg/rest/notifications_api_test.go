package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
	"github.com/owjoel/client-factpack/apps/notif/pkg/rest"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&storage.Notification{})
	return db
}

// setupRouter initializes a router with a test database
func setupRouter() (*gin.Engine, *storage.NotificationStorage) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	
	db := setupTestDB()
	store := &storage.NotificationStorage{DB: db}
	
	rest.InitNotificationAPI(store)
	r.GET("/notifications", rest.GetUserNotifications)
	r.GET("/notifications/client", rest.GetClientNotifications)
	
	return r, store
}

// TestGetUserNotifications tests the user notification endpoint
func TestGetUserNotifications(t *testing.T) {
	r, store := setupRouter()

	t.Run("Missing username", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/notifications", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "username query param is required")
	})

	t.Run("Valid request", func(t *testing.T) {
		// Insert test data
		store.DB.Create(&storage.Notification{
			Username:         "testuser",
			JobID:            "1",
			Status:           "pending",
			NotificationType: "job",
			Type:             "background-job",
		})

		req, _ := http.NewRequest(http.MethodGet, "/notifications?username=testuser", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "pending")
	})

	t.Run("Database error", func(t *testing.T) {
		// Create new router with a DB that will produce errors
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		
		// Create a storage with a closed DB to simulate errors
		db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close() // This will cause DB operations to fail
		
		errorStore := &storage.NotificationStorage{DB: db}
		rest.InitNotificationAPI(errorStore)
		r.GET("/notifications", rest.GetUserNotifications)
		
		req, _ := http.NewRequest(http.MethodGet, "/notifications?username=testuser", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve notifications")
	})
}

// TestGetClientNotifications tests the client notification endpoint
func TestGetClientNotifications(t *testing.T) {
	r, store := setupRouter()

	t.Run("Valid request", func(t *testing.T) {
		// Insert test data
		store.DB.Create(&storage.Notification{
			NotificationType: "client",
			ClientID:         "123",
			ClientName:       "Test Client",
			Priority:         "high",
		})

		req, _ := http.NewRequest(http.MethodGet, "/notifications/client?priority=high", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Client")
	})

	t.Run("Database error", func(t *testing.T) {
		// Create new router with a DB that will produce errors
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		
		// Create a storage with a closed DB to simulate errors
		db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
		sqlDB, _ := db.DB()
		sqlDB.Close() // This will cause DB operations to fail
		
		errorStore := &storage.NotificationStorage{DB: db}
		rest.InitNotificationAPI(errorStore)
		r.GET("/notifications/client", rest.GetClientNotifications)
		
		req, _ := http.NewRequest(http.MethodGet, "/notifications/client", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to retrieve client notifications")
	})
}
