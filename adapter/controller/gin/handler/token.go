package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/api"
	"go-banking-api/entity"
	"go-banking-api/pkg"
	"go-banking-api/pkg/logger"
	"go-banking-api/usecase"
)

type TokenHandler struct {
	tokenUsecase usecase.TokenUsecase
	clock        pkg.Clock
}

func NewTokenHandler(tokenUsecase usecase.TokenUsecase, clock pkg.Clock) *TokenHandler {
	if clock == nil {
		clock = pkg.RealClock{}
	}
	return &TokenHandler{
		tokenUsecase: tokenUsecase,
		clock:        clock,
	}
}

func (t *TokenHandler) PostToken(c *gin.Context) {
	var request presenter.TokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Info(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, "invalid request"))
		return
	}

	token, err := t.tokenUsecase.Refresh(request.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrRefreshTokenRequired):
			logger.Info(err.Error())
			c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, "refresh token is required"))
		case errors.Is(err, usecase.ErrInvalidRefreshToken):
			logger.Info(err.Error())
			c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "invalid refresh token"))
		default:
			logger.Error(err.Error())
			c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, "internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, t.tokenToResponse(token))
}

func (t *TokenHandler) tokenToResponse(token *entity.Token) *presenter.TokenResponse {
	expiresIn := int(token.ExpiresAt.Sub(t.clock.Now()).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	return &presenter.TokenResponse{
		ApiVersion: api.Version,
		Data: presenter.TokenData{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    expiresIn,
		},
	}
}
