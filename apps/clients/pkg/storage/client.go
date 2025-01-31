package storage

import (
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name        string
	Nationality string
}

type ClientInterface interface {
	Create(r *model.CreateClientReq) (model.StatusRes, error)
	Get(clientID string) (*model.GetClientRes)
	Update(r *model.UpdateClientReq) (model.StatusRes, error)
	Delete(r *model.DeleteClientReq) (model.StatusRes, error)
}

type ClientStorage struct {
	*gorm.DB
}

func (c *ClientStorage) Get(clientID string) (*model.GetClientRes) {
	return nil
}

func (c *ClientStorage) Create(r *model.CreateClientReq) (model.StatusRes, error) {
	return model.StatusRes{}, nil
}

func (c *ClientStorage) Update(r *model.UpdateClientReq) (model.StatusRes, error) {
	return model.StatusRes{}, nil
}

func (c *ClientStorage) Delete(r *model.DeleteClientReq) (model.StatusRes, error) {
	return model.StatusRes{}, nil
}