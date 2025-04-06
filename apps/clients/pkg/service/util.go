package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/owjoel/client-factpack/apps/clients/config"
)

func getUsername(ctx context.Context) string {
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return "Unknown"
	}
	return username
}

func TriggerPrefectFlowRun(deploymentID, apiKey string, params map[string]interface{}) error {
	requestBody := map[string]interface{}{
		"parameters": params,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	url := config.PrefectAPIURL + deploymentID + "/create_flow_run"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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