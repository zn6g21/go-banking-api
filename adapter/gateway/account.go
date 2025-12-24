package gateway

import (
	"gorm.io/gorm"

	"go-banking-api/entity"
)

type AccountRepository interface {
	Get(cifNo int) (*entity.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (a *accountRepository) Get(cifNo int) (*entity.Account, error) {
	var account = entity.Account{}
	if err := a.db.Where("cif_no = ?", cifNo).Take(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
