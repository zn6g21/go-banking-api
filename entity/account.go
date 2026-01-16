package entity

import "time"

type AccountStatus string

const (
	AccountStatusActive  AccountStatus = "active"
	AccountStatusClosed  AccountStatus = "closed"
	AccountStatusFrozen  AccountStatus = "frozen"
	AccountStatusDormant AccountStatus = "dormant"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

type Account struct {
	Id            int
	CifNo         int
	Status        AccountStatus
	BranchCode    string
	AccountNumber string
	AccountType   string
	Currency      string
	Balance       int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a *Account) IsActive() bool {
	return a.Status == AccountStatusActive
}
