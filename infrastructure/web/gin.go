package web

import (
	"context"
	"fmt"
	"net/http"

	"go-banking-api/adapter/controller/gin/router"
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

func NewGinServer(host, port string, corsAllowOrigins []string) Server {
	router := router.NewGinRouter(corsAllowOrigins)
	/*if err != nil {
		logger.Error(err.Error(), "host", host, "port", port)
		return nil, err
	} */
	return &GinWebServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", host, port),
			Handler: router,
		},
	}
}
