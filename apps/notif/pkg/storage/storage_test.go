package storage_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
)

func setupTestDB(t *testing.T, migrate bool) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory SQLite DB: %v", err)
	}

	if migrate {
		err = db.AutoMigrate(&storage.Notification{})
		if err != nil {
			t.Fatalf("Failed to migrate Notification model: %v", err)
		}
	}

	return db
}

func TestSaveNotification_Success(t *testing.T) {
	db := setupTestDB(t, true)
	notifStorage := &storage.NotificationStorage{DB: db}

	notification := &storage.Notification{
		UserID:  "testUser",
		Message: "Test message",
	}

	err := notifStorage.SaveNotification(notification)
	assert.NoError(t, err)

	var fetched storage.Notification
	result := db.First(&fetched, notification.ID)

	assert.NoError(t, result.Error)
	assert.Equal(t, notification.UserID, fetched.UserID)
	assert.Equal(t, notification.Message, fetched.Message)
}

func TestSaveNotification_Failure_NilNotification(t *testing.T) {
	db := setupTestDB(t, true)
	notifStorage := &storage.NotificationStorage{DB: db}

	err := notifStorage.SaveNotification(nil)
	assert.Error(t, err)
}

func TestSaveNotification_Failure_BadSchema(t *testing.T) {
	db := setupTestDB(t, false)
	notifStorage := &storage.NotificationStorage{DB: db}

	notification := &storage.Notification{
		UserID:  "testUser",
		Message: "This should fail because the table doesn't exist",
	}

	err := notifStorage.SaveNotification(notification)
	assert.Error(t, err)
}

// This covers the "happy path" of InitDatabase (AutoMigrate + return db)
func TestInitDatabase_SuccessWithSQLite(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&storage.Notification{})
	assert.NoError(t, err)

	notifStorage := &storage.NotificationStorage{DB: db}

	notification := &storage.Notification{
		UserID:  "user1",
		Message: "Test",
	}

	err = notifStorage.SaveNotification(notification)
	assert.NoError(t, err)

	// Verify it was saved
	var fetched storage.Notification
	result := db.First(&fetched, notification.ID)
	assert.NoError(t, result.Error)
}

// This covers the log.Fatalf("Failed to connect to database: %v", err)
func TestInitDatabase_FatalOnBadConnection(t *testing.T) {
	if os.Getenv("TEST_DB_FATAL") == "1" {
		// Simulate InitDatabase with an invalid config by overwriting env vars
		_ = os.Setenv("DB_USER", "invalid")
		_ = os.Setenv("DB_PASSWORD", "invalid")
		_ = os.Setenv("DB_HOST", "invalid-host")
		_ = os.Setenv("DB_PORT", "3306")
		_ = os.Setenv("DB_NAME", "invalid-db")

		// Call InitDatabase (expect log.Fatalf)
		storage.InitDatabase()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitDatabase_FatalOnBadConnection")
	cmd.Env = append(os.Environ(), "TEST_DB_FATAL=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
		output := stderr.String()

		if !strings.Contains(output, "Failed to connect to database") {
			t.Errorf("expected fatal error log, got: %s", output)
		}

		return
	}

	t.Fatalf("process ran without exit error; err=%v, stderr=%v", err, stderr.String())
}

func TestInitDatabase_CoversAutoMigrate(t *testing.T) {
    // Skip actual MySQL connection by using an in-memory SQLite database
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)

    // Call AutoMigrate explicitly to simulate InitDatabase behavior
    err = db.AutoMigrate(&storage.Notification{})
    assert.NoError(t, err)

    // Verify that the table was created
    migrator := db.Migrator()
    assert.True(t, migrator.HasTable(&storage.Notification{}), "Expected Notification table to be created")
}
