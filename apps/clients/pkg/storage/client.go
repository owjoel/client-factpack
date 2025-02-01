package storage

import (
	"fmt"
	"time"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name        string `gorm:"name"`
	Age         uint `gorm:"age"`
	Nationality string `gorm:"nationality"`
	Status string `gorm:"status"`
}

type ClientInterface interface {
	Create(c *model.Client) error
	Get(clientID uint) (*model.Client, error)
	Update(c *model.Client) error
	// Delete(c *model.Client) error
}

type ClientStorage struct {
	*gorm.DB
}

func (s *ClientStorage) Get(clientID uint) (*model.Client, error) {
	var c Client
	if res := s.DB.First(&c, clientID); res.Error != nil {
		return nil, fmt.Errorf("Error retrieving from DB: %w", res.Error)
	}
	return &model.Client{
		Name: c.Name,
		Age: c.Age,
		Nationality: c.Nationality,
		Status: c.Status,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ClientStorage) Create(c *model.Client) error {
	client := Client{
		Name: c.Name,
		Age: c.Age,
		Nationality: c.Nationality,
		Status: c.Status,
	}
	if res := s.DB.Create(&client); res.Error != nil {
		return fmt.Errorf("Error creating client in DB: %w", res.Error)
	}
	return nil
}

func (s *ClientStorage) Update(c *model.Client) error {
	client := &Client{
		Model: gorm.Model{ID: c.ID},
		Name: c.Name,
		Age: c.Age,
		Nationality: c.Nationality,
	}
	if res := s.DB.Model(&client).Updates(client); res.Error != nil {
		return fmt.Errorf("Error updating client in DB: %w", res.Error)
	}
	return nil
}

// func (s *ClientStorage) Delete(c *model.Client) error {
// 	client := &Client{Model: gorm.Model{ID: c.ID}}
// 	if res := s.DB.Model(&client).Update("status", "Inactive"); res.Error != nil {
// 		return fmt.Errorf("Error deactivating client in DB: %w", res.Error)
// 	}
// 	return nil
// }