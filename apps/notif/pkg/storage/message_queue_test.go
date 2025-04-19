package storage

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"os"
	"os/exec"

	"github.com/owjoel/client-factpack/apps/notif/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/notif/config"
	"github.com/stretchr/testify/assert"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestSQLiteDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&Notification{}))
	return db
}

func TestInitMessageQueueIntegration(t *testing.T) {
	// Override RabbitMQ URL
	config.RabbitMQURL = "amqp://guest:guest@localhost:5672/"

	// Set up in-memory SQLite DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&Notification{}))

	// Call InitMessageQueue in background
	go InitMessageQueue(db)
	time.Sleep(500 * time.Millisecond) // Allow listener to start

	// Connect to RabbitMQ
	conn, err := amqp.Dial(config.RabbitMQURL)
	assert.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	assert.NoError(t, err)
	defer ch.Close()

	_, err = ch.QueueDeclare("notifications", true, false, false, false, nil)
	assert.NoError(t, err)

	// Publish a fake message
	msg := NotificationMessage{
		NotificationType: model.NotificationTypeJob,
		Title:            "Integration Test",
		Username:         "testuser",
		ClientName:       []string{"Client A"},
	}
	body, _ := json.Marshal(msg)

	err = ch.PublishWithContext(
		context.Background(),
		"",
		"notifications",
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	assert.NoError(t, err)

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	// Check DB
	var saved Notification
	err = db.First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, "testuser", saved.Username)
	assert.Equal(t, "Integration Test", saved.Title)
	assert.Equal(t, "Client A", saved.ClientName)
}

func TestInitMessageQueue_BadRabbitMQURL(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		db := setupTestSQLiteDB(t)
		config.RabbitMQURL = "invalid-url"
		InitMessageQueue(db)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitMessageQueue_BadRabbitMQURL")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		// expected crash
	} else {
		t.Fatalf("Expected process to exit with error, got: %v", err)
	}
}

func TestInitMessageQueue_FailChannel(t *testing.T) {
	// Setup RabbitMQ with an invalid connection (close immediately)
	conn, _ := amqp.Dial(config.RabbitMQURL)
	conn.Close() // force failure

	// Override URL to hit this closed connection
	config.RabbitMQURL = "amqp://guest:guest@localhost:5672/" // still valid, but closed

	// Capture fatal
	if os.Getenv("BE_CRASHER") == "1" {
		db := setupTestSQLiteDB(t)
		InitMessageQueue(db)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitMessageQueue_FailChannel")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Expected crash, got: %v", err)
	}
}

func TestInitMessageQueue_FailQueueDeclare(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		db := setupTestSQLiteDB(t)

		// Set invalid port to trigger queue declare failure
		config.RabbitMQURL = "amqp://guest:guest@localhost:5673/"

		InitMessageQueue(db) // Should call log.Fatal
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitMessageQueue_FailQueueDeclare")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	err := cmd.Run()

	if err == nil {
		t.Fatalf("Expected queue declare crash, got: %v", err)
	}
}


func TestInitMessageQueue_FailConsume(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		db := setupTestSQLiteDB(t)

		// This port shouldn't have any RabbitMQ server
		config.RabbitMQURL = "amqp://guest:guest@localhost:5674/"

		// This will force failure at or before Consume
		InitMessageQueue(db)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestInitMessageQueue_FailConsume")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Failed to consume message: %v", err)
	}
}

func TestInitMessageQueue_UnmarshalError(t *testing.T) {
	db := setupTestSQLiteDB(t)
	go InitMessageQueue(db)
	time.Sleep(500 * time.Millisecond)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	assert.NoError(t, err)
	defer ch.Close()

	// Broken JSON
	badJSON := []byte(`{"notificationType": job`) // <-- invalid

	err = ch.PublishWithContext(
		context.Background(),
		"",
		"notifications",
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        badJSON,
		},
	)
	assert.NoError(t, err)

	time.Sleep(500 * time.Millisecond)
}

