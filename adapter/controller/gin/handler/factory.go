package handler

type ServerHandler struct {
	*AccountInfoHandler
	*TokenHandler
}

func NewServerHandler(accountInfoHandler *AccountInfoHandler, tokenHandler *TokenHandler) *ServerHandler {
	return &ServerHandler{
		AccountInfoHandler: accountInfoHandler,
		TokenHandler:       tokenHandler,
	}
}
