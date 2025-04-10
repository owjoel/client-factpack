package storage

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/owjoel/client-factpack/apps/notif/pkg/api"
	"github.com/owjoel/client-factpack/apps/notif/pkg/utils"
)

type NotificationMessage struct {
	UserID           string           `json:"userId"`
	NotificationType model.NotificationType `json:"notificationType"`
	Username         string           `json:"username,omitempty"`
	ID               string           `json:"id,omitempty"`
	Status           model.JobStatus  `json:"status,omitempty"`
	Type             model.JobType    `json:"type,omitempty"`
	ClientName       string           `json:"clientName,omitempty"`
	Priority         model.Priority   `json:"priority,omitempty"`
}

// Initialize RabbitMQ listener
func InitMessageQueue(db *gorm.DB) {
	conn, err := amqp.Dial(config.RabbitMQURL)
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

	store := &NotificationStorage{db}

	utils.Logger.Info("Listening for messages...")

	go func() {
		for msg := range msgs {
			var notification NotificationMessage
			err := json.Unmarshal(msg.Body, &notification)
			if err != nil {
				utils.Logger.Error("Error parsing message:", err)
				continue
			}

			// Store in DB
			store.SaveNotification(&Notification{
				UserID:           notification.UserID,
				NotificationType: string(notification.NotificationType),
				Username:         notification.Username,
				ID:               notification.ID,
				Status:           string(notification.Status),
				Type:             string(notification.Type),
				ClientName:       notification.ClientName,
				Priority:         string(notification.Priority),
			})

			// Forward to WebSocket
			msgBytes, _ := json.Marshal(notification)
			api.SendNotification(notification.UserID, string(msgBytes))
		}
	}()
}