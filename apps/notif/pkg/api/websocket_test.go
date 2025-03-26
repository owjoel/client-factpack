package api_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
)

func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(api.HandleWebSocketConnections))
}

func TestHandleWebSocketConnections_MissingUserID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	wsURL := "ws" + server.URL[4:] + "/ws"

	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)
	assert.NotNil(t, err)
	assert.Nil(t, conn)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleWebSocketConnections_SuccessfulConnection(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	u := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}
	q := u.Query()
	q.Set("userId", "testUser")
	u.RawQuery = q.Encode()

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer ws.Close()

	time.Sleep(100 * time.Millisecond)

	err = ws.WriteMessage(websocket.TextMessage, []byte("Hello"))
	assert.NoError(t, err)
}

func TestSendNotification_UserConnected(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	u := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}
	q := u.Query()
	q.Set("userId", "testUser")
	u.RawQuery = q.Encode()

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer ws.Close()

	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		_, msg, err := ws.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, "Test message", string(msg))
	}()

	api.SendNotification("testUser", "Test message")

	wg.Wait()
}

func TestSendNotification_UserNotConnected(t *testing.T) {
	// No user connected
	// Assuming SendNotification does not return an error for non-connected users
	api.SendNotification("ghostUser", "No one is here!")
}

func TestHandleWebSocketConnections_UpgradeError(t *testing.T) {
	// Create a test HTTP server with a handler that triggers the upgrade error.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intentionally call the handler with a bad HTTP method (non-GET), which WebSocket upgrader rejects.
		r.Method = http.MethodPost // WebSocket upgrader requires GET
		api.HandleWebSocketConnections(w, r)
	}))
	defer server.Close()

	// Perform a POST request instead of WebSocket dial (simulate a client not upgrading)
	resp, err := http.Post(server.URL+"?userId=testUser", "application/json", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check that we get an error, and the status is 400/500 depending on the server setup.
	assert.NotEqual(t, http.StatusSwitchingProtocols, resp.StatusCode)
}

func TestSendNotification_MultipleUsers(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	u1 := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}
	q1 := u1.Query()
	q1.Set("userId", "user1")
	u1.RawQuery = q1.Encode()

	u2 := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}
	q2 := u2.Query()
	q2.Set("userId", "user2")
	u2.RawQuery = q2.Encode()

	ws1, _, err := websocket.DefaultDialer.Dial(u1.String(), nil)
	assert.NoError(t, err)
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(u2.String(), nil)
	assert.NoError(t, err)
	defer ws2.Close()

	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		_, msg, err := ws1.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, "Hello user1", string(msg))
	}()

	go func() {
		defer wg.Done()

		_, msg, err := ws2.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, "Hello user2", string(msg))
	}()

	api.SendNotification("user1", "Hello user1")
	api.SendNotification("user2", "Hello user2")

	wg.Wait()
}

func TestSendNotification_DisconnectedUser(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	u := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}
	q := u.Query()
	q.Set("userId", "testUser")
	u.RawQuery = q.Encode()

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Close the connection
	ws.Close()

	// Try sending a notification to the disconnected user
	api.SendNotification("testUser", "This should not be received")
}

func TestSendNotification_ErrorSendingMessage(t *testing.T) {
    // Set up the test WebSocket server
    server := setupTestServer()
    defer server.Close()

    // Create a WebSocket connection
    u := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/ws"}
    q := u.Query()
    q.Set("userId", "testUser")
    u.RawQuery = q.Encode()

    ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    assert.NoError(t, err)

    // Close the WebSocket connection to simulate a broken pipe
    ws.Close()
    time.Sleep(100 * time.Millisecond) // Ensure the connection is fully closed

    // Send a notification to the disconnected user
    api.SendNotification("testUser", "This will fail to send")

    // The logger should log an error when `WriteMessage` fails
    // You can verify this manually by checking the logs or mock the logger if needed.
}

func TestStartWebSocketServer(t *testing.T) {
	// This calls StartWebSocketServer, which registers the handler and logs the startup message.
	api.StartWebSocketServer()

	// Optionally, make a test request to confirm handler registration
	req := httptest.NewRequest(http.MethodGet, "/ws?userId=testUser", nil)
	w := httptest.NewRecorder()

	api.HandleWebSocketConnections(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode) // Because no proper upgrade attempt, it returns error
}
