package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/owjoel/client-factpack/apps/clients/config"
	errorx "github.com/owjoel/client-factpack/apps/clients/pkg/api/errors"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientService struct {
	clientRepository repository.ClientRepository
	jobService       JobServiceInterface
	logService       LogServiceInterface
	prefectFlowRunner   PrefectFlowRunnerInterface
}

type ClientServiceInterface interface {
	GetClient(ctx context.Context, clientID string) (*model.Client, error)
	GetAllClients(ctx context.Context, query *model.GetClientsQuery) (total int, clients []model.Client, err error)
	CreateClientByName(ctx context.Context, req *model.CreateClientByNameReq) (string, error)
	UpdateClient(ctx context.Context, clientID string, changes []model.SimpleChanges) error
	MatchClient(ctx context.Context, req *model.MatchClientReq, clientID string) (string, error)
}

func NewClientService(clientRepository repository.ClientRepository, jobService JobServiceInterface, logService LogServiceInterface, prefectFlowRunner PrefectFlowRunnerInterface) *ClientService {
	return &ClientService{clientRepository: clientRepository, jobService: jobService, logService: logService, prefectFlowRunner: prefectFlowRunner}
}

func (s *ClientService) GetClient(ctx context.Context, clientID string) (*model.Client, error) {
	c, err := s.clientRepository.GetOne(ctx, clientID)
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) || errors.Is(err, errorx.ErrDependencyFailed) || errors.Is(err, errorx.ErrInvalidInput) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: error getting client", errorx.ErrInternal)
	}

	username := GetUsername(ctx)
	_, err = s.logService.CreateLog(ctx, &model.Log{
		ClientID:  clientID,
		Actor:     username,
		Operation: model.OperationGet,
		Details:   fmt.Sprintf("User %s viewed client profile with id %s", username, clientID),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("error creating log: %v", err) // don't return error since it's not critical
	}

	return c, nil
}

func (s *ClientService) GetAllClients(ctx context.Context, query *model.GetClientsQuery) (total int, clients []model.Client, err error) {
	clients, err = s.clientRepository.GetAll(ctx, query)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return 0, nil, err
		}
		return 0, nil, fmt.Errorf("%w: error getting clients", errorx.ErrInternal)
	}

	total, err = s.clientRepository.Count(ctx, query)

	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return 0, nil, err
		}
		return 0, nil, fmt.Errorf("%w: error getting clients", errorx.ErrInternal)
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
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return "", err
		}
		return "", fmt.Errorf("%w: error creating job", errorx.ErrInternal)
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
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return "", err
		}
		return "", fmt.Errorf("%w: error creating client", errorx.ErrInternal)
	}

	// trigger prefect workflow with job id
	err = s.prefectFlowRunner.Trigger(config.PrefectScrapeFlowID, map[string]interface{}{
		"job_id":    id,
		"target":    req.Name,
		"client_id": clientId,
		"username":  GetUsername(ctx),
	})
	if err != nil {
		return "", fmt.Errorf("%w: error triggering prefect workflow", errorx.ErrInternal)
	}

	username := GetUsername(ctx)
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
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return err
		}
		return fmt.Errorf("%w: error creating job", errorx.ErrInternal)
	}

	clientName, err := s.clientRepository.GetClientNameByID(ctx, clientID)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) || errors.Is(err, errorx.ErrNotFound) || errors.Is(err, errorx.ErrInvalidInput) {
			return err
		}
		return fmt.Errorf("%w: error getting client name", errorx.ErrInternal)
	}

	err = s.prefectFlowRunner.Trigger(
		config.PrefectScrapeFlowID,
		map[string]interface{}{
			"job_id":    id,
			"target":    clientName,
			"client_id": clientID,
			"username":  GetUsername(ctx),
		},
	)
	if err != nil {
		return fmt.Errorf("%w: error triggering prefect workflow", errorx.ErrInternal)
	}

	username := GetUsername(ctx)
	_, err = s.logService.CreateLog(ctx, &model.Log{
		ClientID:  clientID,
		Actor:     username,
		Operation: model.OperationScrape,
		Details:   fmt.Sprintf("User %s rescrapped client profile with job id %s", username, id),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("error creating log: %v", err) // don't return error since it's not critical
	}

	return nil
}

func (s *ClientService) UpdateClient(ctx context.Context, clientID string, changes []model.SimpleChanges) error {
	client, err := s.clientRepository.GetOne(ctx, clientID)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) || errors.Is(err, errorx.ErrNotFound) || errors.Is(err, errorx.ErrInvalidInput) {
			return err
		}
		return fmt.Errorf("%w: error getting client", errorx.ErrInternal)
	}

	if client == nil {
		return errorx.ErrNotFound
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
		return errorx.ErrInvalidInput
	}

	if err := s.clientRepository.Update(ctx, clientID, update); err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) || errors.Is(err, errorx.ErrInvalidInput) {
			return err
		}
		return fmt.Errorf("%w: error updating client", errorx.ErrInternal)
	}

	username := GetUsername(ctx)
	_, err = s.logService.CreateLog(ctx, &model.Log{
		ClientID:  clientID,
		Actor:     username,
		Operation: model.OperationUpdate,
		Details:   fmt.Sprintf("User %s updated client profile with id %s", username, clientID),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("error creating log: %v", err) // don't return error since it's not critical
	}

	return nil
}

func (s *ClientService) MatchClient(ctx context.Context, req *model.MatchClientReq, clientID string) (string, error) {
	job := &model.Job{
		Type:      model.Match,
		Status:    model.JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Logs: []model.JobLog{
			{
				Message:   "Job [MATCH] created and submitted to Prefect",
				Timestamp: time.Now(),
			},
		},
	}

	id, err := s.jobService.CreateJob(ctx, job)
	if err != nil {
		if errors.Is(err, errorx.ErrDependencyFailed) {
			return "", err
		}
		return "", fmt.Errorf("%w: error creating job", errorx.ErrInternal)
	}

	err = s.prefectFlowRunner.Trigger(
		config.PrefectMatchFlowID,
		map[string]interface{}{
			"job_id":     id,
			"file_name":  req.FileName,
			"file_bytes": req.FileBytes,
			"target_id": clientID,
			"username":  GetUsername(ctx),
		},
	)
	if err != nil {
		return "", fmt.Errorf("%w: error triggering prefect workflow", errorx.ErrInternal)
	}

	return id, nil
}
