package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-banking-api/entity"
	"go-banking-api/pkg"
)

func TestCustomer(t *testing.T) {
	now := pkg.Str2time("2025-12-02")
	birthDate := pkg.Str2time("1990-01-01")
	customer := entity.Customer{
		CifNo:      1,
		NameKana:   "Taro Tanaka",
		NameKanji:  "田中 太郎",
		BirthDate:  birthDate,
		Prefecture: "Tokyo",
		City:       "Bunkyo",
		Town:       "Kouraku",
		Street:     "1-1-1",
		Building:   "CivicCenter",
		Room:       "101",
		Email:      "tarou.tanaka@example.com",
		Phone:      "09012345678",
		CreatedAt:  now,
	}

	assert.Equal(t, 1, customer.CifNo)
	assert.Equal(t, "Taro Tanaka", customer.NameKana)
	assert.Equal(t, "田中 太郎", customer.NameKanji)
	assert.Equal(t, birthDate, customer.BirthDate)
	assert.Equal(t, "Tokyo", customer.Prefecture)
	assert.Equal(t, "Bunkyo", customer.City)
	assert.Equal(t, "Kouraku", customer.Town)
	assert.Equal(t, "1-1-1", customer.Street)
	assert.Equal(t, "CivicCenter", customer.Building)
	assert.Equal(t, "101", customer.Room)
	assert.Equal(t, "tarou.tanaka@example.com", customer.Email)
	assert.Equal(t, "09012345678", customer.Phone)
	assert.Equal(t, now, customer.CreatedAt)
}
