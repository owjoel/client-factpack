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
	NotificationType string `gorm:"column:notification_type"`
	Username         string `gorm:"column:username"`
	JobID               string `gorm:"column:job_id"`
	Status           string `gorm:"column:job_status"`
	Type             string `gorm:"column:job_type"`
	ClientID       string `gorm:"column:client_id"`
	ClientName       string `gorm:"column:client_name"`
	Priority         string `gorm:"column:priority"`
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
		log.Fatalf(config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&Notification{})
	return db
}

// Saves notification to the database
func (s *NotificationStorage) SaveNotification(n *Notification) error {
	return s.Create(n).Error
}

func (s *NotificationStorage) GetNotificationsByUser(userID string) ([]Notification, error) {
	var notifications []Notification
	err := s.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error
	return notifications, err
}
