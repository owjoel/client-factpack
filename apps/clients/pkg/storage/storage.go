package storage

import (

	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
	"gorm.io/gorm"
)

type storage struct {
	Client ClientInterface
}

var (
	db *storage
)

func Init() {
	clientStorage := InitMySQL()
	db = &storage{Client: clientStorage}
}

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
}

// username:password@protocol(address)/dbname?param=value

func GetInstance() *storage {
	return db
}
