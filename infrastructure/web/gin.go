package web

import (
	"context"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"go-banking-api/adapter/controller/gin/router"
	"go-banking-api/pkg/logger"
)

type GinWebServer struct {
	server *http.Server
}

func (g *GinWebServer) Start() error {
	return g.server.ListenAndServe()
}

func (g *GinWebServer) Shutdown(ctx context.Context) error {
	return g.server.Shutdown(ctx)
}

func NewGinServer(host, port string, corsAllowOrigins []string, db *gorm.DB) (Server, error) {
	router, err := router.NewGinRouter(db, corsAllowOrigins)
	if err != nil {
		logger.Error(err.Error(), "host", host, "port", port)
		return nil, err
	}
	return &GinWebServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", host, port),
			Handler: router,
		},
	}, err
}
