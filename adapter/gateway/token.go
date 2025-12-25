package gateway

import (
	"go-banking-api/entity"

	"gorm.io/gorm"
)

type TokenRepository interface {
	Get(token string) (*entity.Token, error)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (t *tokenRepository) Get(tokenVal string) (*entity.Token, error) {
	var token = entity.Token{}
	if err := t.db.Where("access_token = ?", tokenVal).Take(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}
