package service_test

// import (
// 	"testing"
// 	"bytes"
// 	"os"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/owjoel/client-factpack/apps/notif/pkg/service"
// )

// type MockNotifier struct {
// 	Called         bool
// 	CapturedUserID string
// 	CapturedMsg    string
// }

// func (m *MockNotifier) SendNotification(userID, message string) {
// 	m.Called = true
// 	m.CapturedUserID = userID
// 	m.CapturedMsg = message
// }

// func TestSendNotification_CallsNotifier(t *testing.T) {
// 	mockNotifier := &MockNotifier{}
// 	notificationService := &service.NotificationService{
// 		Notifier: mockNotifier,
// 	}

// 	notificationService.SendNotification("testUser", "hello world")

// 	assert.True(t, mockNotifier.Called)
// 	assert.Equal(t, "testUser", mockNotifier.CapturedUserID)
// 	assert.Equal(t, "hello world", mockNotifier.CapturedMsg)
// }

// func TestAPINotifier_SendNotification(t *testing.T) {
// 	// Set up a pipe to capture stdout
// 	r, w, _ := os.Pipe()
// 	stdout := os.Stdout
// 	os.Stdout = w

// 	// Run the test code
// 	notifier := &service.APINotifier{}
// 	notifier.SendNotification("testUser", "Hello from APINotifier")

// 	// Restore stdout and close writer
// 	w.Close()
// 	os.Stdout = stdout

// 	// Read from pipe
// 	var buf bytes.Buffer
// 	_, _ = buf.ReadFrom(r)
// 	output := buf.String()

// 	// Assert fmt.Println was called
// 	assert.Contains(t, output, "Sending notification to user: testUser")
// };