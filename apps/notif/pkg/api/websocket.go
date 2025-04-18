package api

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/owjoel/client-factpack/apps/notif/pkg/utils"
)

var (
	clients   = make(map[string]*websocket.Conn) // Map userID -> WebSocket connection
	clientsMu sync.Mutex                         // Mutex for safe access
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// Handle WebSocket connections
func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("username")
	if userID == "" {
		utils.Logger.Warn("WebSocket connection attempt without username")
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Logger.Error("WebSocket upgrade error: ", err)
		return
	}
	defer conn.Close()

	clientsMu.Lock()
	clients[userID] = conn
	clientsMu.Unlock()

	utils.Logger.Infof("User %s connected to WebSocket", userID)

	// Handle disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			utils.Logger.Infof("User %s disconnected", userID)
			clientsMu.Lock()
			delete(clients, userID)
			clientsMu.Unlock()
			break
		}
	}
}

// Send notification to a specific user
func SendNotification(userID, message, notificationType string) {
	clientsMu.Lock()
	conn, exists := clients[userID]
	utils.Logger.Info("Clients", clients, userID)
	clientsMu.Unlock()

	if notificationType == "client" {
		for _, conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				utils.Logger.Info("Error sending message:", err)
			}
		}
	} else if exists {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			utils.Logger.Info("Error sending message:", err)
		}
	} else {
		utils.Logger.Info("User not connected:", userID)
	}
}

// Start WebSocket server
func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleWebSocketConnections)
	utils.Logger.Info("WebSocket server running...")
}
