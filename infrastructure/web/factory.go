package web

import (
	"context"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

func NewServer() Server {
	config := NewConfigWeb()
	return NewGinServer(config.Host, config.Port, config.CorsAllowOrigins)
}
