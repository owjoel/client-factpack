package model

type NotificationType string
type JobStatus string
type JobType string
type Priority string

const (
	NotificationTypeJob    NotificationType = "job"
	NotificationTypeClient NotificationType = "client"

	JobStatusCompleted  JobStatus = "completed"
	JobStatusPending    JobStatus = "pending"
	JobStatusFailed     JobStatus = "failed"
	JobStatusProcessing JobStatus = "processing"

	JobTypeScrape JobType = "scrape"
	JobTypeMatch  JobType = "match"

	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

type Notification struct {
	NotificationType NotificationType `json:"notificationType"`
	Username         string           `json:"username,omitempty"`    // job
	JobID               string           `json:"id,omitempty"`          // job
	Status           JobStatus        `json:"status,omitempty"`      // job
	Type             JobType          `json:"type,omitempty"`        // job
	ClientID       string		   `json:"clientId,omitempty"`    // client
	ClientName       string           `json:"clientName,omitempty"`  // client
	Priority         Priority         `json:"priority,omitempty"`    // client
}