package handler

import "github.com/gin-gonic/gin"

type APIHandler struct {
	accountInfo *AccountInfoHandler
	token       *TokenHandler
}

func NewAPIHandler(accountInfo *AccountInfoHandler, token *TokenHandler) *APIHandler {
	return &APIHandler{
		accountInfo: accountInfo,
		token:       token,
	}
}

func (h *APIHandler) GetAccountInformation(c *gin.Context) {
	h.accountInfo.GetAccountInformation(c)
}

func (h *APIHandler) GetTransactionList(c *gin.Context) {
	h.accountInfo.GetTransactionList(c)
}

func (h *APIHandler) PostToken(c *gin.Context) {
	h.token.PostToken(c)
}
