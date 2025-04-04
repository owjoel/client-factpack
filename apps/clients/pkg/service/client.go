package service

import (
	"context"
	"fmt"
	"time"

	"github.com/owjoel/client-factpack/apps/clients/config"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientService struct {
	clientRepository repository.ClientRepository
	jobService       JobServiceInterface
	logService       LogServiceInterface
}

type ClientServiceInterface interface {
	GetClient(ctx context.Context, clientID string) (*model.Client, error)
	GetAllClients(ctx context.Context, query *model.GetClientsQuery) (total int, clients []model.Client, err error)
	CreateClientByName(ctx context.Context, req *model.CreateClientByNameReq) error
	UpdateClient(ctx context.Context, clientID string, changes []model.SimpleChanges) error
}

func NewClientService(clientRepository repository.ClientRepository, jobService JobServiceInterface, logService LogServiceInterface) *ClientService {
	return &ClientService{clientRepository: clientRepository, jobService: jobService, logService: logService}
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

	total, err = s.clientRepository.Count(ctx, query)

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
		Logs: []model.JobLog{
			{
				Message:   "Job [CREATE] created and submitted to Prefect",
				Timestamp: time.Now(),
			},
		},
	}
	id, err := s.jobService.CreateJob(ctx, job)
	if err != nil {
		return "", fmt.Errorf("Error creating job: %w", err)
	}

	// create client profile
	client := &model.Client{
		Data: bson.D{
			{
				Key: "profile", Value: bson.D{
					{Key: "names", Value: bson.A{req.Name}},
				},
			},
		},
		Metadata: model.ClientMetadata{
			Scraped:   false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	clientId, err := s.clientRepository.Create(ctx, client)
	if err != nil {
		return "", fmt.Errorf("error creating client profile: %w", err)
	}

	// trigger prefect workflow with job id
	go TriggerScrapeFlowRun(config.PrefectAPIURL, config.PrefectScrapeFlowID, config.PrefectAPIKey, map[string]interface{}{
		"job_id":    id,
		"target":    req.Name,
		"client_id": clientId,
	})

	username := getUsername(ctx)
	s.logService.CreateLog(ctx, &model.Log{
		ClientID:  clientId,
		Actor:     username,
		Operation: model.OperationCreateAndScrape,
		Details:   fmt.Sprintf("User %s created a new client profile with job id %s", username, id),
		Timestamp: time.Now(),
	})

	return id, nil
}

func (s *ClientService) RescrapeClient(ctx context.Context, clientID string) error {
	job := &model.Job{
		Type:      model.Scrape,
		Status:    model.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Logs: []model.JobLog{
			{
				Message:   "Job [RESCRAPE] created and submitted to Prefect",
				Timestamp: time.Now(),
			},
		},
	}

	id, err := s.jobService.CreateJob(ctx, job)
	if err != nil {
		return fmt.Errorf("error creating job: %w", err)
	}

	clientName, err := s.clientRepository.GetClientNameByID(ctx, clientID)
	if err != nil {
		return fmt.Errorf("error getting client name: %w", err)
	}

	go TriggerScrapeFlowRun(
		config.PrefectAPIURL,
		config.PrefectScrapeFlowID,
		config.PrefectAPIKey,
		map[string]interface{}{
			"job_id":    id,
			"target":    clientName,
			"client_id": clientID,
		},
	)

	return nil
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
