package entity

import (
	"go-banking-api/pkg"
	"strings"
	"time"
)

type Token struct {
	AccessToken string
	RefreshToken string
	Scopes      string // "read:account_and_transactions write:transfer" のようなスペース区切り
	ExpiresAt   time.Time
	CifNo       int
}

func (t *Token) IsExpired(clock pkg.Clock) bool {
	return clock.Now().After(t.ExpiresAt)
}

func (t *Token) HasScope(targetScope string) bool {
	scopes := strings.Split(t.Scopes, " ")
	for _, s := range scopes {
		if s == targetScope {
			return true
		}
	}
	return false
}
