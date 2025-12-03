package gateway

import (
	"gorm.io/gorm"

	"go-banking-api/entity"
)

type CustomerRepository interface {
	Get(cifNo int) (*entity.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (c *customerRepository) Get(cifNo int) (*entity.Customer, error) {
	var customer = entity.Customer{}
	if err := c.db.Take(&customer, cifNo).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}
