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
	ID           bson.ObjectID   `bson:"_id,omitempty" json:"id" swaggerignore:"true"`
	Type         JobType         `bson:"type" json:"type"`
	Status       JobStatus       `bson:"status" json:"status"`
	CreatedAt    time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time       `bson:"updatedAt" json:"updatedAt"`
	Input        bson.M          `bson:"input" json:"input"`
	ScrapeResult bson.ObjectID   `bson:"scrapeResult" json:"scrapeResult"`
	MatchResults []bson.ObjectID `bson:"matchResults" json:"matchResults"`
}

type GetJobsQuery struct {
	Status   JobStatus `bson:"status" json:"status"`
	Page     int       `bson:"page" json:"page"`
	PageSize int       `bson:"pageSize" json:"pageSize"`
}

type GetJobsResponse struct {
	Jobs []Job `bson:"jobs" json:"jobs"`
}
