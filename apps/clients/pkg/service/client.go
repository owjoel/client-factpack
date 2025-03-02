package service

import (
	"context"
	"fmt"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"github.com/owjoel/client-factpack/apps/clients/pkg/storage"
)

type ClientService struct {
	storage storage.ClientInterface
}

func NewClientService(storage storage.ClientInterface) *ClientService {
	return &ClientService{storage: storage}
}

func (s *ClientService) GetClient(ctx context.Context, clientID string) (*model.Client, error) {
	c, err := storage.GetInstance().Client.Get(ctx, clientID)
	if err != nil {
		return &model.Client{}, fmt.Errorf("Error retrieving client profile: %w", err)
	}
	return c, nil
}

func (s *ClientService) GetAllClients(ctx context.Context) ([]model.Client, error) {
	clients, err := storage.GetInstance().Client.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving all client records: %w", err)
	}
	return clients, nil
}

func (s *ClientService) CreateClient(ctx context.Context, client *model.Client) error {
	err := storage.GetInstance().Client.Create(ctx, client)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	return nil
}

// func (s *ClientService) UpdateClient(ctx context.Context, r model.UpdateClientReq) (*model.StatusRes, error) {
// 	client := &model.Client{
// 		Name:        r.Name,
// 		Age:         r.Age,
// 		Nationality: r.Nationality,
// 	}
// 	if err := storage.GetInstance().Client.Update(ctx, client); err != nil {
// 		return nil, fmt.Errorf("Error updating client profile: %w", err)
// 	}
// 	return &model.StatusRes{Status: "Success"}, nil
// }
