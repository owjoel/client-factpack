package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Log struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id" swaggerignore:"true"`
	ClientID  string `bson:"clientID" json:"clientID"`
	Actor     string        `bson:"actor" json:"actor"`
	Operation Operation     `bson:"operation" json:"operation"`
	Details   string        `bson:"details" json:"details"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
}

type Operation string

const (
	OperationCreate Operation = "create"
	OperationUpdate Operation = "update"
	OperationDelete Operation = "delete"
	OperationScrape Operation = "scrape"
	OperationMatch  Operation = "match"
	OperationCreateAndScrape Operation = "create & scrape"
)

type GetLogsQuery struct {
	ClientID  string    `bson:"clientID" json:"clientID"`
	Operation Operation `bson:"operation" json:"operation"`
	Actor     string    `bson:"actor" json:"actor"` // username of the actor
	From      time.Time `bson:"from" json:"from"`
	To        time.Time `bson:"to" json:"to"`
	Page      int       `bson:"page" json:"page"`
	PageSize  int       `bson:"pageSize" json:"pageSize"`
}

type GetLogsResponse struct {
	Total int    `json:"total"`
	Logs  []Log  `json:"logs"`
}


type GetLogResponse struct {
	Log *Log `json:"log"`
}

type CreateLogResponse struct {
	ID string `json:"id"`
}
