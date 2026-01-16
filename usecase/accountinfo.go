package usecase

import (
	"errors"

	"go-banking-api/adapter/gateway"
	"go-banking-api/entity"

	"gorm.io/gorm"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountInactive = errors.New("account is not active")
)

type AccountInfo struct {
	NameKana      string
	NameKanji     string
	Status        entity.AccountStatus
	BranchCode    string
	AccountNumber string
	AccountType   string
	Currency      string
	Balance       int64
}

type AccountInfoUsecase interface {
	Get(cifNo int) (*AccountInfo, error)
}

type accountInfoUsecase struct {
	customerRepository gateway.CustomerRepository
	accountRepository  gateway.AccountRepository
}

func NewAccountInfoUsecase(
	customerRepository gateway.CustomerRepository,
	accountRepository gateway.AccountRepository,
) *accountInfoUsecase {
	return &accountInfoUsecase{
		customerRepository: customerRepository,
		accountRepository:  accountRepository,
	}
}

func (a *accountInfoUsecase) Get(cifNo int) (*AccountInfo, error) {
	customer, err := a.customerRepository.Get(cifNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	account, err := a.accountRepository.Get(cifNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	if !account.IsActive() {
		return nil, ErrAccountInactive
	}
	return &AccountInfo{
		NameKana:      customer.NameKana,
		NameKanji:     customer.NameKanji,
		Status:        account.Status,
		BranchCode:    account.BranchCode,
		AccountNumber: account.AccountNumber,
		AccountType:   account.AccountType,
		Currency:      account.Currency,
		Balance:       account.Balance,
	}, nil
}
