package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Client struct {
	ID       bson.ObjectID  `bson:"_id,omitempty" json:"id" swaggerignore:"true"`
	Data     bson.D         `bson:"data" json:"data"`
	Metadata ClientMetadata `bson:"metadata" json:"metadata"`
}

type ClientMetadata struct {
	Scraped   bool      `bson:"scraped" json:"scraped"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	Sources   []string  `bson:"sources" json:"sources"`
}

// Request-response models

type GetClientsQuery struct {
	Name     string `form:"name"`
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
}

type GetClientsResponse struct {
	Total int      `json:"total"`
	Data  []Client `json:"data"`
}

type StatusRes struct {
	Status string `json:"status"`
}

type GetClientRes struct {
	Name        string
	Age         uint
	Nationality string
}

type CreateClientReq struct {
	Data bson.M `json:"data"`
}

type CreateClientByNameReq struct {
	Name string `json:"name"`
}

type CreateClientByNameRes struct {
	JobID string `json:"jobId"`
}

type UpdateClientReq struct {
	Changes []SimpleChanges `json:"changes"`
}

type SimpleChanges struct {
	Path  string `json:"path"`
	Old   any    `json:"old"`
	New   any    `json:"new"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// // Client contains all information for a particular client
// type Client struct {
// 	// gorm.Model
// 	ID          bson.ObjectID `bson:"_id,omitempty" swaggerignore:"true"`
// 	Profile     Profile       `json:"profile"`
// 	Investments []Investment  `json:"investments"`
// 	Associates  []Associate   `json:"associates"`
// 	Metadata    Metadata      `json:"metadata"`

// 	Status string `gorm:"status"`
// }

// // Profile contains basic personal information about the client
// type Profile struct {
// 	Name             string        `json:"name" example:"john doe"`
// 	Age              uint          `json:"age" example:"55"`
// 	Nationality      string        `json:"nationality" example:"chinese"`
// 	CurrentResidence Residence     `bson:"currentResidence" json:"currentResidence"`
// 	NetWorth         NetWorth      `json:"netWorth"`
// 	Industries       []string      `json:"industries"`
// 	Occupations      []string      `json:"occupations"`
// 	Socials          []SocialMedia `json:"socials"`
// 	Contact          Contact       `json:"contact"`
// }

// // Residence contains the details of a client's residence
// type Residence struct {
// 	City    string `json:"city"`
// 	Country string `json:"country"`
// }

// // NetWorth contains information about the client's current net worth
// type NetWorth struct {
// 	EstimatedValue uint      `bson:"estimatedValue" json:"estimatedValue"`
// 	Currency       string    `json:"currency"`
// 	Source         string    `json:"source"`
// 	Timestamp      time.Time `bson:"timestamp" json:"timestamp"`
// }

// // SocialMedia contains information about a client's current
// type SocialMedia struct {
// 	Platform string `json:"platform"`
// 	Username string `json:"username"`
// }

// // Contact contains a client's contact information
// type Contact struct {
// 	WorkAddress string `bson:"workAddress" json:"workAddress"`
// 	Phone       string `json:"phone"`
// }

// // Investment contains information about a client's investment
// type Investment struct {
// 	Name     string          `json:"name"`
// 	Type     string          `json:"type"`
// 	Value    InvestmentValue `json:"value"`
// 	Date     time.Time       `json:"date"`
// 	Industry string          `json:"industry"`
// 	Status   string          `json:"status"`
// 	Source   string          `json:"source"`
// }

// // InvestmentValue contains the value and currency for a particular investment
// type InvestmentValue struct {
// 	Value    uint   `json:"value"`
// 	Currency string `json:"currency"`
// }

// // Associate contains information about a client's known associate
// type Associate struct {
// 	Name                string   `json:"name"`
// 	Relationship        string   `json:"relationship"`
// 	AssociatedCompanies []string `bson:"associatedCompanies" json:"associatedCompanies"`
// }

// // Metadata contains general information about the client's profile in the app.
// type Metadata struct {
// 	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
// 	Sources   []string  `json:"sources"`
// }
