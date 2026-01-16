package web

import (
	"context"

	"gorm.io/gorm"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func NewServer(db *gorm.DB) (Server, error) {
	config := NewConfigWeb()
	return NewGinServer(config.Host, config.Port, config.CorsAllowOrigins, db)
}
