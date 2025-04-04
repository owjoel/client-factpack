package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/owjoel/client-factpack/apps/notif/pkg/service"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api/model"

)

// MockNotifier implements service.Notifier for testing
type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) SendNotification(userID, message string) {
	m.Called(userID, message)
}

// SetupRouter manually defines a test POST route for integration testing
func setupRouter(svc *service.NotificationService) *gin.Engine {
	r := gin.Default()
	r.POST("/api/v1/notify", func(c *gin.Context) {
		var req struct {
			UserID  string `json:"userId"`
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		svc.SendNotification(req.UserID, req.Message)
		c.JSON(http.StatusAccepted, gin.H{"status": "sent"})
	})
	return r
}

func TestSendNotificationAPI(t *testing.T) {
	mockNotifier := new(MockNotifier)
	svc := &service.NotificationService{Notifier: mockNotifier}

	router := setupRouter(svc)

	payload := map[string]string{
		"userId":  "123",
		"message": "Test message",
	}
	body, _ := json.Marshal(payload)

	// Expect the mock to be called with these values
	mockNotifier.On("SendNotification", "123", "Test message").Return()

	req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
	mockNotifier.AssertExpectations(t)
}

func TestWebSocketNotificationDelivery(t *testing.T) {
	// Start test server using the WebSocket handler
	s := httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
	defer s.Close()

	// Build WebSocket URL
	wsURL := "ws" + s.URL[4:] + "/?userId=testuser"

	// Connect to WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Give server time to register the client
	time.Sleep(100 * time.Millisecond)

	// Send a notification to the connected user
	api.SendNotification("testuser", "Hello from test!")

	// Read message from WebSocket
	_, msg, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "Hello from test!", string(msg))
}

func TestSendNotificationAPI_MissingFields(t *testing.T) {
	mockNotifier := new(MockNotifier)
	svc := &service.NotificationService{Notifier: mockNotifier}
	router := setupRouter(svc)

	// Missing "message" field
	payload := map[string]string{
		"userId": "123",
	}
	body, _ := json.Marshal(payload)

	mockNotifier.On("SendNotification", "123", "").Return()
	req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
}

func TestSendNotificationAPI_InvalidJSON(t *testing.T) {
	mockNotifier := new(MockNotifier)
	svc := &service.NotificationService{Notifier: mockNotifier}
	router := setupRouter(svc)

	// Invalid JSON payload (missing closing brace)
	body := []byte(`{"userId": "123", "message": "oops"`)

	req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestWebSocketNotificationDelivery_MissingUserID(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
	defer s.Close()

	wsURL := "ws" + s.URL[4:] // no userId

	// Try connecting without userId
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

	// Expect handshake failure
	assert.Error(t, err)
	assert.Nil(t, ws)ose()
}

func TestRabbitMQToWebSocketIntegration(t *testing.T) {
	// Start the WebSocket server
	wsServer := httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
	defer wsServer.Close()

	wsURL := "ws" + wsServer.URL[4:] + "/?userId=testuser"

	// Connect WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Start the real RabbitMQ consumer (using real notifier.SendNotification)
	go func() {
		storage.InitMessageQueue()
	}()

	time.Sleep(200 * time.Millisecond) // Give consumer time to start

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)
	defer conn.Close()

	ch, err := conn.Channel()
	assert.NoError(t, err)
	defer ch.Close()

	_, err = ch.QueueDeclare("notifications", false, false, false, false, nil)
	assert.NoError(t, err)

	// Prepare the message
	msg := model.Notification{
		UserID:  "testuser",
		Message: "Hello via RabbitMQ!",
	}
	body, _ := json.Marshal(msg)

	// Publish message to queue
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

	// Read message from WebSocket
	_, resp, err := ws.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, "Hello via RabbitMQ!", string(resp))
}