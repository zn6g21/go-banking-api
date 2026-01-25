package handler

import "sync"

var (
	serverHandler *ServerHandler
	once          sync.Once
)

type ServerHandler struct {
	*AccountInfoHandler
	*TokenHandler
}

func NewHandler() *ServerHandler {
	once.Do(func() {
		serverHandler = &ServerHandler{}
	})
	return serverHandler
}

func (h *ServerHandler) Register(i interface{}) *ServerHandler {
	switch interfaceType := i.(type) {
	case *AccountInfoHandler:
		serverHandler.AccountInfoHandler = interfaceType
	case *TokenHandler:
		serverHandler.TokenHandler = interfaceType
	}
	return serverHandler
}
