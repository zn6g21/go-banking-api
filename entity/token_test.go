package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-banking-api/entity"
)

func TestToken(t *testing.T) {
	now := time.Now()
	token := entity.Token{
		AccessToken: "test-token",
		Scopes:      "read:account_and_transactions write:transfer",
		ExpiresAt:   now,
	}

	assert.Equal(t, "test-token", token.AccessToken)
	assert.Equal(t, "read:account_and_transactions write:transfer", token.Scopes)
	assert.Equal(t, now, token.ExpiresAt)
}

func TestIsExpired(t *testing.T) {
	token := entity.Token{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(t, token.IsExpired())

	token.ExpiresAt = time.Now().Add(1 * time.Hour)
	assert.False(t, token.IsExpired())
}

func TestHasScope(t *testing.T) {
	token := entity.Token{
		Scopes: "read:account_and_transactions write:transfer",
	}

	assert.True(t, token.HasScope("read:account_and_transactions"))
	assert.True(t, token.HasScope("write:transfer"))
	assert.False(t, token.HasScope("read"))

	token.Scopes = ""
	assert.False(t, token.HasScope("read:account_and_transactions"))
}
