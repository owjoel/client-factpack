package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&Notification{}))
	return db
}

func TestSaveNotification(t *testing.T) {
	db := setupTestDB(t)
	store := &NotificationStorage{DB: db}

	notification := &Notification{
		NotificationType: "job",
		Title:            "Test Notification",
		Source:           "system",
		Username:         "user1",
		JobID:            "job-001",
		Status:           "completed",
		Type:             "match",
		ClientID:         "client-001",
		ClientName:       "Client A;Client B",
		Priority:         "high",
	}

	err := store.SaveNotification(notification)
	assert.NoError(t, err)

	var result Notification
	err = db.First(&result).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Notification", result.Title)
}

func TestGetNotificationsByUser(t *testing.T) {
	db := setupTestDB(t)
	store := &NotificationStorage{DB: db}

	db.Create(&Notification{
		NotificationType: "job",
		Username:         "john",
		JobID:            "job-001",
		Status:           "completed",
		Type:             "scrape",
	})

	notif, err := store.GetNotificationsByUser("john", "completed", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, notif, 1)
	assert.Equal(t, "job-001", notif[0].JobID)
}

func TestGetClientNotifications(t *testing.T) {
	db := setupTestDB(t)
	store := &NotificationStorage{DB: db}

	db.Create(&Notification{
		NotificationType: "client",
		Title:            "Client Alert",
		Source:           "sourceX",
		ClientID:         "client-xyz",
		ClientName:       "Client A;Client B",
		Priority:         "medium",
	})

	notif, err := store.GetClientNotifications("Client A", "medium", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, notif, 1)
	assert.Equal(t, "Client Alert", notif[0].Title)
	assert.Equal(t, []string{"Client A", "Client B"}, notif[0].ClientName)
}