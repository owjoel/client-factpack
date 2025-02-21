package storage

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
	"github.com/owjoel/client-factpack/apps/notif/pkg/utils"
)

type NotificationMessage struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

// Initialize RabbitMQ listener
func InitMessageQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		utils.Logger.Fatal("Failed to connect to RabbitMQ:", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		utils.Logger.Fatal("Failed to connect to channel:", err)
		return
	}

	q, err := ch.QueueDeclare("notifications", false, false, false, false, nil)
	if err != nil {
		utils.Logger.Fatal("Failed to declare queue:", err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		utils.Logger.Fatal("Failed to consume message:", err)
		return
	}

	utils.Logger.Info("Listening for messages...")

	go func() {
		for msg := range msgs {
			var notification NotificationMessage
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				utils.Logger.Fatal("Error parsing message:", err)
				continue
			}

			// Send notification to WebSocket client
			utils.Logger.Infof("Message received for User %s: %s", notification.UserID, notification.Message)
			api.SendNotification(notification.UserID, notification.Message)
		}
	}()
}
