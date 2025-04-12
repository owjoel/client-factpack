package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Log struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id" swaggerignore:"true"`
	ClientID  string        `bson:"clientId" json:"clientId"`
	Actor     string        `bson:"actor" json:"actor"`
	Operation Operation     `bson:"operation" json:"operation"`
	Details   string        `bson:"details" json:"details"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
}

type Operation string

const (
	OperationGet             Operation = "view"
	OperationCreate          Operation = "create"
	OperationUpdate          Operation = "update"
	OperationDelete          Operation = "delete"
	OperationScrape          Operation = "scrape"
	OperationMatch           Operation = "match"
	OperationCreateAndScrape Operation = "create & scrape"
)

type GetLogsQuery struct {
	ClientID  string    `bson:"clientId" json:"clientId" form:"clientId"`
	Operation Operation `bson:"operation" json:"operation" form:"operation"`
	Actor     string    `bson:"actor" json:"actor" form:"actor"` // username of the actor
	From      time.Time `bson:"from" json:"from" form:"from"`
	To        time.Time `bson:"to" json:"to" form:"to"`
	Page      int       `bson:"page" json:"page" form:"page"`
	PageSize  int       `bson:"pageSize" json:"pageSize" form:"pageSize"`
}

type GetLogsResponse struct {
	Total int   `json:"total"`
	Logs  []Log `json:"logs"`
}

type GetLogResponse struct {
	Log *Log `json:"log"`
}

type CreateLogResponse struct {
	ID string `json:"id"`
}
