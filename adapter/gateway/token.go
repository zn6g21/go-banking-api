package gateway

import (
	"time"

	"go-banking-api/entity"

	"gorm.io/gorm"
)

type TokenRepository interface {
	Get(token string) (*entity.Token, error)
	UpdateByRefreshToken(refreshToken string, accessToken string, newRefreshToken string, expiresAt time.Time) error
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

func (t *tokenRepository) UpdateByRefreshToken(refreshToken string, accessToken string, newRefreshToken string, expiresAt time.Time) error {
	result := t.db.Model(&entity.Token{}).Where("refresh_token = ?", refreshToken).Updates(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"expires_at":    expiresAt,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
