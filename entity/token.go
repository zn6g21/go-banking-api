package entity

import (
	"strings"
	"time"
)

type Token struct {
	AccessToken string
	Scopes      string // "read:account_and_transactions write:transfer" のようなスペース区切り
	ExpiresAt   time.Time
	CifNo       int
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
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
