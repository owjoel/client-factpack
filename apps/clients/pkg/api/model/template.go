package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty" swaggerignore:"true"`
	Version   int                    `bson:"version" json:"version"`
	Schema    map[string]interface{} `bson:"schema" json:"schema"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
}
