package service

import (
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

func (s *ClientService) CreateClient(r *model.CreateClientReq) (*model.StatusRes, error) {
	client := &model.Client{
		Name:        r.Name,
		Age:         r.Age,
		Nationality: r.Nationality,
	}

	if err := s.storage.Create(client); err != nil {
		return &model.StatusRes{}, fmt.Errorf("Error creating client profile: %w", err)
	}
	return &model.StatusRes{Status: "Success"}, nil
}

func (s *ClientService) GetClient(clientID uint) (*model.GetClientRes, error) {
	c, err := s.storage.Get(clientID)
	if err != nil {
		return &model.GetClientRes{}, fmt.Errorf("Error retrieving client profile: %w", err)
	}
	return &model.GetClientRes{
		Name:        c.Name,
		Age:         c.Age,
		Nationality: c.Nationality,
	}, nil
}

func (s *ClientService) UpdateClient(r model.UpdateClientReq) (*model.StatusRes, error) {
	client := &model.Client{
		Name:        r.Name,
		Age:         r.Age,
		Nationality: r.Nationality,
	}
	if err := s.storage.Update(client); err != nil {
		return nil, fmt.Errorf("Error updating client profile: %w", err)
	}
	return &model.StatusRes{Status: "Success"}, nil
}

// func (s *ClientService) DeleteClient(r model.DeleteClientReq) (*model.StatusRes, error) {
// 	client := &model.Client{
// 		ID: r.ID,
// 	}

// 	if err := s.storage.Delete(client); err != nil {
// 		return nil, fmt.Errorf("Error deactivating client profile: %w", err)
// 	}
// 	return &model.StatusRes{Status: "Success"}, nil
// }
