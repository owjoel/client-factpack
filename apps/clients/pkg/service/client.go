package service

import (
	"context"
	"fmt"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClientService struct {
	clientRepository repository.ClientRepository
	jobRepository repository.JobRepository
}

type ClientServiceInterface interface {
	GetClient(ctx context.Context, clientID string) (*model.Client, error)
	GetAllClients(ctx context.Context, query *model.GetClientsQuery) (total int, clients []model.Client, err error)
	CreateClientByName(ctx context.Context, req *model.CreateClientByNameReq) error
	UpdateClient(ctx context.Context, clientID string, data bson.D) error
}

func NewClientService(clientRepository repository.ClientRepository, jobRepository repository.JobRepository) *ClientService {
	return &ClientService{clientRepository: clientRepository, jobRepository: jobRepository}
}

func (s *ClientService) GetClient(ctx context.Context, clientID string) (*model.Client, error) {
	c, err := s.clientRepository.GetOne(ctx, clientID)
	if err != nil {
		return &model.Client{}, fmt.Errorf("Error retrieving client profile: %w", err)
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

func (s *ClientService) CreateClientByName(ctx context.Context, req *model.CreateClientByNameReq) error {
	// create job
	// trigger prefect workflow with job id
	// return job id
	return nil
}

func (s *ClientService) UpdateClient(ctx context.Context, clientID string, data bson.D) error {
	client, err := s.clientRepository.GetOne(ctx, clientID)
	if err != nil {
		return fmt.Errorf("Error updating client: %w", err)
	}

	if client == nil {
		return fmt.Errorf("Client not found")
	}

	err = s.clientRepository.Update(ctx, clientID, data)
	
	return nil
}
