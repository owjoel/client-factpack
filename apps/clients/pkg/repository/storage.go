package repository

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	// "gorm.io/gorm"
// )

// type storage struct {
// 	Client ClientRepository
// }

// var (
// 	db *storage
// )

// func Init() {
// 	clientStorage := InitMongo()
// 	db = &storage{Client: clientStorage}
// }

// // type ClientRepository interface {
// // 	Create(ctx context.Context, c *model.Client) error
// // 	GetOne(ctx context.Context, clientID string) (*model.Client, error)
// // 	GetAll(ctx context.Context) ([]model.Client, error)
// // 	Update(ctx context.Context, c *model.Client) error
// // }

// func NewClientRepository() ClientRepository {
// 	return InitMongo()
// }

// func GetInstance() *storage {
// 	return db
// }

// // Adding this function to allow injection for testing.
// func SetInstanceClient(client ClientRepository) {
// 	if db == nil {
// 		db = &storage{}
// 	}
// 	db.Client = client
// }
