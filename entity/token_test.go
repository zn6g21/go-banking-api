package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-banking-api/entity"
	"go-banking-api/pkg"
)

func TestToken(t *testing.T) {
	now := time.Now()
	token := entity.Token{
		AccessToken:  "test-token",
		RefreshToken: "refresh-token-1",
		Scopes:       "read:account_and_transactions write:transfer",
		ExpiresAt:    now,
		CifNo:        1,
		ClientID:     "client-1",
	}

	assert.Equal(t, "test-token", token.AccessToken)
	assert.Equal(t, "refresh-token-1", token.RefreshToken)
	assert.Equal(t, "read:account_and_transactions write:transfer", token.Scopes)
	assert.Equal(t, now, token.ExpiresAt)
	assert.Equal(t, 1, token.CifNo)
	assert.Equal(t, "client-1", token.ClientID)
}

func TestIsExpired(t *testing.T) {
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}

	token := entity.Token{
		ExpiresAt: fixedNow.Add(-1 * time.Hour),
	}
	assert.True(t, token.IsExpired(clock))

	token.ExpiresAt = fixedNow.Add(1 * time.Hour)
	assert.False(t, token.IsExpired(clock))
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
