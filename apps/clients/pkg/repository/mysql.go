package repository

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/owjoel/client-factpack/apps/clients/config"
// 	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// type SQLStorage struct {
// 	*gorm.DB
// }

// func GetDSN() string {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
// 		config.DBUser,
// 		config.DBPassword,
// 		config.DBHost,
// 		config.DBPort,
// 		config.DBName,
// 	)
// 	return dsn
// }

// func InitMySQL() *SQLStorage {
// 	_db, err := gorm.Open(mysql.Open(GetDSN()), &gorm.Config{})
// 	if err != nil {
// 		panic("Failed to connect to database.")
// 	}
// 	log.Printf("Connected to DB")
// 	if err := _db.AutoMigrate(&Client{}); err != nil {
// 		panic("Failed to migrate resource model.")
// 	}

// 	return &SQLStorage{_db}
// }

// func (s *SQLStorage) Get(clientID uint) (*model.Client, error) {
// 	var c Client
// 	if res := s.DB.First(&c, clientID); res.Error != nil {
// 		return nil, fmt.Errorf("Error retrieving from DB: %w", res.Error)
// 	}
// 	return &model.Client{
// 		Name:        c.Name,
// 		Age:         c.Age,
// 		Nationality: c.Nationality,
// 		Status:      c.Status,
// 		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
// 		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
// 	}, nil
// }

// func (s *SQLStorage) Create(c *model.Client) error {
// 	client := Client{
// 		Name:        c.Name,
// 		Age:         c.Age,
// 		Nationality: c.Nationality,
// 		Status:      c.Status,
// 	}
// 	if res := s.DB.Create(&client); res.Error != nil {
// 		return fmt.Errorf("Error creating client in DB: %w", res.Error)
// 	}
// 	return nil
// }

// func (s *SQLStorage) Update(c *model.Client) error {
// 	client := &Client{
// 		Model:       gorm.Model{ID: c.ID},
// 		Name:        c.Name,
// 		Age:         c.Age,
// 		Nationality: c.Nationality,
// 	}
// 	if res := s.DB.Model(&client).Updates(client); res.Error != nil {
// 		return fmt.Errorf("Error updating client in DB: %w", res.Error)
// 	}
// 	return nil
// }
