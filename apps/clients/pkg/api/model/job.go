package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

type JobType string

const (
	Scrape JobType = "scrape"
	Match  JobType = "match"
)

type Job struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id" swaggerignore:"true"`
	PrefectFlowID string        `bson:"prefectFlowID" json:"prefectFlowID"`
	Type          JobType       `bson:"type" json:"type"`
	Status        JobStatus     `bson:"status" json:"status"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
	Input         bson.M        `bson:"input" json:"input"`
	ScrapeResult  bson.ObjectID `bson:"scrapeResult" json:"scrapeResult"`
	MatchResults  []MatchResult `bson:"matchResults" json:"matchResults"`
	Logs          []JobLog      `bson:"logs" json:"logs"`
}

type JobLog struct {
	Message   string        `bson:"message" json:"message"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
}

type MatchResult struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id" swaggerignore:"true"`
	ConfidenceScore      float64       `bson:"confidenceScore" json:"confidenceScore"`
}

type GetJobsQuery struct {
	Status   JobStatus `bson:"status" json:"status" form:"status"`
	Page     int       `bson:"page" json:"page" form:"page"`
	PageSize int       `bson:"pageSize" json:"pageSize" form:"pageSize"`
}

type GetJobsResponse struct {
	Total int   `bson:"total" json:"total"`
	Jobs  []Job `bson:"jobs" json:"jobs"`
}
