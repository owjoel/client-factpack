package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
	"github.com/owjoel/client-factpack/apps/notif/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MockNotifier implements service.Notifier
type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) SendNotification(userID, message, notifType string) {
	m.Called(userID, message, notifType)
}

func setupRouter(svc *MockNotifier) *gin.Engine {
	r := gin.Default()
	r.POST("/api/v1/notify", func(c *gin.Context) {
		var req struct {
			UserID  string `json:"userId"`
			Message string `json:"message"`
			Type    string `json:"type"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		svc.SendNotification(req.UserID, req.Message, req.Type)
		c.JSON(http.StatusAccepted, gin.H{"status": "sent"})
	})
	return r
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	err = db.AutoMigrate(&storage.Notification{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	return db
}

func TestSendNotificationAPI(t *testing.T) {
	mockNotifier := new(MockNotifier)
	svc := mockNotifier
	router := setupRouter(svc)

	payload := map[string]string{
		"userId":  "123",
		"message": "Test message",
		"type":    "job",
	}
	body, _ := json.Marshal(payload)
	mockNotifier.On("SendNotification", "123", "Test message", "job").Return()

	req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
	mockNotifier.AssertExpectations(t)
}

func TestSendNotificationAPI_InvalidJSON(t *testing.T) {
	mockNotifier := new(MockNotifier)
	svc := mockNotifier
	router := setupRouter(svc)

	body := []byte(`{"userId": "123", "message": "oops"`)

	req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWebSocketNotificationDelivery(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
	defer s.Close()

	// Include required userId
	wsURL := "ws" + s.URL[4:] + "/?username=testuser"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to establish WebSocket connection: %v", err)
	}
	defer ws.Close()

	// Send a notification
	api.SendNotification("testuser", "Hello from test!", "job")

	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "Hello from test!", string(msg))
}

func TestWebSocketNotificationDelivery_MissingUserID(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
	defer s.Close()

	wsURL := "ws" + s.URL[4:]
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Error(t, err)
	assert.Nil(t, ws)
}

func TestRabbitMQToWebSocketIntegration(t *testing.T) {
	// Load environment configs
	config.Load()

	// Start WebSocket test server
	wsServer := httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
	defer wsServer.Close()

	wsURL := "ws" + wsServer.URL[4:] + "/?username=testuser"

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Start the RabbitMQ consumer with test DB
	go func() {
		db := setupTestDB(t)
		storage.InitMessageQueue(db)
	}()

	time.Sleep(200 * time.Millisecond) // Give consumer time to start

	// Connect to RabbitMQ
	conn, err := amqp.Dial(config.RabbitMQURL)
	assert.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	assert.NoError(t, err)
	defer ch.Close()

	_, err = ch.QueueDeclare("notifications", true, false, false, false, nil)
	assert.NoError(t, err)

	// Prepare message body
	msg := storage.NotificationMessage{
		NotificationType: model.NotificationTypeJob,
		Title:            "Hello via RabbitMQ",
		Source:           "integration-test",
		Username:         "testuser",
		JobID:            "job-999",
		Status:           model.JobStatusCompleted,
		Type:             model.JobTypeMatch,
		ClientID:         "client-1",
		ClientName:       []string{"Client A"},
		Priority:         model.PriorityHigh,
	}
	body, _ := json.Marshal(msg)

	// Publish to RabbitMQ
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

	// Expect message on WebSocket
	_, resp, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "Hello via RabbitMQ")
}