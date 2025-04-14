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
	JobID            string `gorm:"column:job_id"`
	Status           string `gorm:"column:status"`
	Type             string `gorm:"column:type"`
	ClientID         string `gorm:"column:client_id"`
	ClientName       string `gorm:"column:client_name"`
	Priority         string `gorm:"column:priority"`
}

type JobNotification struct {
	Username string `json:"username"`
	JobID    string `json:"jobId"`
	Status   string `json:"status"`
	Type     string `json:"type"`
}

type ClientNotification struct {
	ClientID   string `json:"clientId"`   // comes from Notification.UserID
	ClientName string `json:"clientName"` // Notification.ClientName
	Priority   string `json:"priority"`   // Notification.Priority
	JobID      string `json:"jobId"`      // Notification.ID
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

func (s *NotificationStorage) GetNotificationsByUser(username string) ([]JobNotification, error) {
	var result []JobNotification
	err := s.Model(&Notification{}).
		Select("username, job_id AS job_id, status, type").
		Where("username = ? AND notification_type = ?", username, "job").
		Order("created_at DESC").
		Scan(&result).Error

	return result, err
}

func (s *NotificationStorage) GetClientNotifications() ([]ClientNotification, error) {
	var result []ClientNotification
	err := s.Model(&Notification{}).
		Select("client_id AS client_id, client_name, priority, job_id AS job_id").
		Where("notification_type = ?", "client").
		Order("created_at DESC").
		Scan(&result).Error

	return result, err
}
