package model

import "go.mongodb.org/mongo-driver/v2/bson"

type GetArticlesReq struct {
	ID []string `json:"id"`
}

type GetArticlesRes struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	ID      bson.ObjectID `bson:"_id" json:"id"`
	Source  string        `bson:"source" json:"source"`
	Title   string        `bson:"title" json:"title"`
	URL     string        `bson:"url" json:"url"`
	Summary string        `bson:"summary" json:"summary"`
}

type Sentiment struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}
