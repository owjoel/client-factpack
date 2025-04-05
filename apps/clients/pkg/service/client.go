package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/owjoel/client-factpack/apps/clients/config"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientService struct {
	clientRepository repository.ClientRepository
	jobService       JobServiceInterface
}

type ClientServiceInterface interface {
	GetClient(ctx context.Context, clientID string) (*model.Client, error)
	GetAllClients(ctx context.Context, query *model.GetClientsQuery) (total int, clients []model.Client, err error)
	CreateClientByName(ctx context.Context, req *model.CreateClientByNameReq) error
	UpdateClient(ctx context.Context, clientID string, changes []model.SimpleChanges) error
}

func NewClientService(clientRepository repository.ClientRepository, jobService JobServiceInterface) *ClientService {
	return &ClientService{clientRepository: clientRepository, jobService: jobService}
}

func (s *ClientService) GetClient(ctx context.Context, clientID string) (*model.Client, error) {
	c, err := s.clientRepository.GetOne(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving client profile: %w", err)
	}
	return c, nil
}

func (s *ClientService) GetAllClients(ctx context.Context, query *model.GetClientsQuery) (total int, clients []model.Client, err error) {
	clients, err = s.clientRepository.GetAll(ctx, query)
	if err != nil {
		return 0, nil, fmt.Errorf("Error retrieving all client records: %w", err)
	}

	total, err = s.clientRepository.Count(ctx)

	if err != nil {
		return 0, nil, fmt.Errorf("Error retrieving total client records: %w", err)
	}

	return total, clients, nil
}

func (s *ClientService) CreateClientByName(ctx context.Context, req *model.CreateClientByNameReq) (string, error) {
	job := &model.Job{
		Type:      model.Scrape,
		Status:    model.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Logs: []model.Log{
			{
				Message:   "Job created and submitted to Prefect",
				Timestamp: time.Now(),
			},
		},
	}
	id, err := s.jobService.CreateJob(ctx, job)
	if err != nil {
		return "", fmt.Errorf("Error creating job: %w", err)
	}

	// trigger prefect workflow with job id
	// TODO: temporary http request to prefect workflow
	go func() {
		requestBody := map[string]interface{}{
			"parameters": map[string]interface{}{
				"job_id":  id,
				"target": req.Name,
			},
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("Error marshalling request body: %v", err)
			return
		}

		url := config.PrefectAPIURL + config.PrefectScrapeFlowID + "/create_flow_run"

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+config.PrefectAPIKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Printf("HTTP POST failed: %v", err)
			return
		}
		defer resp.Body.Close()
	}()

	return id, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, clientID string, changes []model.SimpleChanges) error {
	client, err := s.clientRepository.GetOne(ctx, clientID)
	if err != nil {
		return fmt.Errorf("error retrieving client: %w", err)
	}

	if client == nil {
		return fmt.Errorf("client not found")
	}

	update := bson.D{}
	for _, change := range changes {
		if change.Path == "" {
			continue
		}
		// Prefix with "data." to target fields inside the data object
		key := "data." + change.Path
		update = append(update, bson.E{Key: key, Value: change.New})
	}

	if len(update) == 0 {
		return fmt.Errorf("no valid changes provided")
	}

	if err := s.clientRepository.Update(ctx, clientID, update); err != nil {
		return fmt.Errorf("error updating client: %w", err)
	}

	return nil
}

