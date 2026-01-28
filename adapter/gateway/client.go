package gateway

import (
	"go-banking-api/entity"

	"gorm.io/gorm"
)

type ClientRepository interface {
	Get(clientID string) (*entity.Client, error)
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (c *clientRepository) Get(clientID string) (*entity.Client, error) {
	var client entity.Client
	if err := c.db.Where("client_id = ?", clientID).Take(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
