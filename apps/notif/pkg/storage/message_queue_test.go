package storage_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/owjoel/client-factpack/apps/notif/pkg/storage"
	"github.com/owjoel/client-factpack/apps/notif/pkg/utils"
)

func TestInitMessageQueue_InfoLogs(t *testing.T) {
	var buf bytes.Buffer

	// Redirect logger output to buffer
	originalOutput := utils.Logger.Out
	utils.Logger.Out = &buf
	defer func() { utils.Logger.Out = originalOutput }()

	// Start the queue listener in a goroutine
	go func() {
		storage.InitMessageQueue()
	}()

	// Give time for connection & setup logs
	time.Sleep(1 * time.Second)

	output := buf.String()

	if !strings.Contains(output, "Listening for messages") {
		t.Error("Expected log 'Listening for messages...' not found")
	}

	// Set up RabbitMQ connection to publish a message
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatal("Failed to connect to RabbitMQ. Is it running locally?")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// Declare queue (should match the one in InitMessageQueue)
	q, err := ch.QueueDeclare("notifications", false, false, false, false, nil)
	if err != nil {
		t.Fatal("Failed to declare queue:", err)
	}

	// Publish a test message
	testMsg := map[string]string{
		"userId":  "testUser",
		"message": "Hello World",
	}
	msgBody, err := json.Marshal(testMsg)
	if err != nil {
		t.Fatal("Failed to marshal test message:", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
		},
	)
	if err != nil {
		t.Fatal("Failed to publish test message:", err)
	}

	// Wait for message processing and log writing
	time.Sleep(2 * time.Second)

	output = buf.String()

	if !strings.Contains(output, "Message received for User testUser") {
		t.Error("Expected log 'Message received for User testUser' not found")
	}
}

