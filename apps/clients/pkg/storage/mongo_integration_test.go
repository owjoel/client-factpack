package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/owjoel/client-factpack/apps/clients/pkg/api/model"
)

func TestClientMongoIntegration(t *testing.T) {
	
	db := InitMongo()
	ctx := context.TODO()

	// Create test client
	mockClient := &model.Client{
		Profile: model.Profile{
			Name:        "Integration Test User",
			Age:         30,
			Nationality: "Testland",
			Contact: model.Contact{
				Phone:       "1234567890",
				WorkAddress: "123 Test St",
			},
		},
		Status: "Active",
		Metadata: model.Metadata{
			UpdatedAt: time.Now(),
			Sources:   []string{"test"},
		},
	}

	// Create
	err := db.Create(ctx, mockClient)
	assert.NoError(t, err, "Insert should succeed")

	// GetAll
	allClients, err := db.GetAll(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, allClients)

	// Get last one
	last := allClients[len(allClients)-1]
	fetched, err := db.Get(ctx, last.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, mockClient.Profile.Name, fetched.Profile.Name)

	//Clean up
	_, err = db.clientCollection.DeleteMany(ctx, map[string]interface{}{
		"profile.name": mockClient.Profile.Name,
	})
	assert.NoError(t, err, "Cleanup failed")
}
