package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PrefectFlowRunnerInterface interface {
	Trigger(deploymentID string, params map[string]interface{}) error
}


type PrefectFlowRunner struct {
	APIURL string
	APIKey string
	Client *http.Client
}

func NewPrefectFlowRunner(apiURL string, apiKey string, client *http.Client) *PrefectFlowRunner {
	return &PrefectFlowRunner{
		APIURL: apiURL,
		APIKey: apiKey,
		Client: client,
	}
}

func (r *PrefectFlowRunner) Trigger(deploymentID string, params map[string]interface{}) error {
	requestBody := map[string]interface{}{
		"parameters": params,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	url := r.APIURL + deploymentID + "/create_flow_run"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("received non-2xx response: %s", resp.Status)
	}

	log.Printf("Triggered Prefect flow run. Status: %s", resp.Status)
	return nil
}
