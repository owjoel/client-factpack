package storage

import (
	"fmt"
	"log"

	"github.com/owjoel/client-factpack/apps/notif/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID  string `gorm:"user_id"`
	Message string `gorm:"message"`
}

type NotificationStorage struct {
	*gorm.DB
}

// Connects to MySQL database
func InitDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&Notification{})
	return db
}

// Saves notification to the database
func (s *NotificationStorage) SaveNotification(n *Notification) error {
	return s.Create(n).Error
}
