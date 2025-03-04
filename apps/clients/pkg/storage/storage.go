package storage

import (
	"context"

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	// "gorm.io/gorm"
)

type storage struct {
	Client ClientInterface
}

var (
	db *storage
)

func Init() {
	clientStorage := InitMongo()
	db = &storage{Client: clientStorage}
}

type ClientInterface interface {
	Create(ctx context.Context, c *model.Client) error
	Get(ctx context.Context, clientID string) (*model.Client, error)
	GetAll(ctx context.Context) ([]model.Client, error)
	Update(ctx context.Context, c *model.Client) error
}

// username:password@protocol(address)/dbname?param=value

func GetInstance() *storage {
	return db
}
