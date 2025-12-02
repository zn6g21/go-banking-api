package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-banking-api/entity"
	"go-banking-api/pkg"
)

func TestAccount(t *testing.T) {
	now := pkg.Str2time("2025-12-02")
	account := entity.Account{
		Id:            1,
		CifNo:         1,
		Status:        entity.AccountStatusActive,
		BranchCode:    "001",
		AccountNumber: "1234567",
		AccountType:   "1",
		Currency:      "JPY",
		Balance:       int64(10000),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.CifNo)
	assert.Equal(t, entity.AccountStatusActive, account.Status)
	assert.Equal(t, "001", account.BranchCode)
	assert.Equal(t, "1234567", account.AccountNumber)
	assert.Equal(t, "1", account.AccountType)
	assert.Equal(t, "JPY", account.Currency)
	assert.Equal(t, int64(10000), account.Balance)
	assert.Equal(t, now, account.CreatedAt)
	assert.Equal(t, now, account.UpdatedAt)
}

func TestIsActive(t *testing.T) {
	account := entity.Account{
		Status: entity.AccountStatusActive,
	}
	assert.True(t, account.IsActive())

	account.Status = entity.AccountStatusClosed
	assert.False(t, account.IsActive())
}
