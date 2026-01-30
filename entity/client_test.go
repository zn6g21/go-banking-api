package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-banking-api/entity"
)

func TestClient(t *testing.T) {
	client := entity.Client{
		ClientID:     "client-123",
		ClientSecret: "secret-xyz",
		ClientName:   "Test Client",
		Scope:        "read:account_and_transactions",
	}

	assert.Equal(t, "client-123", client.ClientID)
	assert.Equal(t, "secret-xyz", client.ClientSecret)
	assert.Equal(t, "Test Client", client.ClientName)
	assert.Equal(t, "read:account_and_transactions", client.Scope)
}
