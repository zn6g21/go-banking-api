package gateway

import (
	"gorm.io/gorm"

	"go-banking-api/entity"
)

type AccountRepository interface {
	Get(Id int) (*entity.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (a *accountRepository) Get(id int) (*entity.Account, error) {
	var account = entity.Account{}
	if err := a.db.Take(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
