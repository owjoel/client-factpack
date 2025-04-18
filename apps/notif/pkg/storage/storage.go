package storage

import (
	"fmt"
	"log"
	"strings"
	"github.com/owjoel/client-factpack/apps/notif/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	NotificationType string `gorm:"column:notification_type"`
	Title            string `gorm:"column:title"`
	Source           string `gorm:"column:source"`
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
	Title      string `json:"title"`
	Source     string `json:"source"`
	ClientID   string `json:"clientId"`   // comes from Notification.UserID
	ClientName []string `json:"clientName"` // Notification.ClientName
	Priority   string `json:"priority"`   // Notification.Priority

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

func (s *NotificationStorage) GetNotificationsByUser(username string, status string, page, pageSize int) ([]JobNotification, error) {
	var result []JobNotification
	query := s.Model(&Notification{}).
		Select("username, job_id, status, type").
		Where("username = ? AND notification_type = ?", username, "job")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&result).Error

	return result, err
}

func (s *NotificationStorage) GetClientNotifications(name, priority string, page, pageSize int) ([]ClientNotification, error) {
	
	type rawClientNotification struct {
		Title      string
		Source     string
		ClientID   string
		ClientName string
		Priority   string
	}
	var rawResult []rawClientNotification
	query := s.Model(&Notification{}).
		Select("client_id, client_name, priority, title, source").
		Where("notification_type = ?", "client")

	if name != "" {
		query = query.Where("LOWER(client_name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&rawResult).Error

		if err != nil {
			return nil, err
		}
	
		// Step 2: Convert to final []ClientNotification
		var result []ClientNotification
		for _, r := range rawResult {
			names := []string{}
			if r.ClientName != "" {
				names = strings.Split(r.ClientName, ";")
			}
	
			result = append(result, ClientNotification{
				Title:      r.Title,
				Source:     r.Source,
				ClientID:   r.ClientID,
				ClientName: names,
				Priority:   r.Priority,
			})
		}
	return result, nil
}