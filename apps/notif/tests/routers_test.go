package web_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/owjoel/client-factpack/apps/notif/pkg/web"
)

func TestHealthCheckRoute_LiveServer(t *testing.T) {
	// Start the router in a goroutine
	go func() {
		web.InitRouter() // runs router.Run(":8081") and blocks
	}()

	// Give the server time to start (not ideal, but works)
	time.Sleep(1 * time.Second)

	// Send a real HTTP request
	resp, err := http.Get("http://localhost:8081/api/v1/health")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) == "" || !containsOK(body) {
		t.Errorf("Expected body to contain 'OK', got %s", string(body))
	}
}

func containsOK(body []byte) bool {
	return string(body) == `{"status":"OK"}` || string(body) == "{\"status\":\"OK\"}"
}
