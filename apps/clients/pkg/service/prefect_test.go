package service_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owjoel/client-factpack/apps/clients/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestPrefectFlowRunner_Trigger_Success(t *testing.T) {
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		err = json.Unmarshal(body, &receivedBody)
		assert.NoError(t, err)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	runner := service.NewPrefectFlowRunner(server.URL+"/", "test-api-key", server.Client())

	params := map[string]interface{}{
		"job_id":    "123",
		"client_id": "abc",
	}

	err := runner.Trigger("deployment-xyz", params)

	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		"parameters": map[string]interface{}{
			"job_id":    "123",
			"client_id": "abc",
		},
	}, receivedBody)
}

func TestPrefectFlowRunner_Trigger_HttpError(t *testing.T) {
	runner := service.NewPrefectFlowRunner("http://invalid-host/", "key", &http.Client{})

	err := runner.Trigger("deployment-xyz", map[string]interface{}{"job_id": "fail"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP request failed")
}

func TestPrefectFlowRunner_Trigger_BadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	runner := service.NewPrefectFlowRunner(server.URL+"/", "key", server.Client())

	err := runner.Trigger("deployment-xyz", map[string]interface{}{"job_id": "fail"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "received non-2xx response")
}
