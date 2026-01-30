package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/api"
	"go-banking-api/entity"

	"go-banking-api/pkg"
	"go-banking-api/pkg/logger"
	"go-banking-api/usecase"
)

type TokenHandler struct {
	tokenUsecase  usecase.TokenUsecase
	clientUsecase usecase.ClientUsecase
	clock         pkg.Clock
}

func NewTokenHandler(tokenUsecase usecase.TokenUsecase, clientUsecase usecase.ClientUsecase, clock pkg.Clock) *TokenHandler {
	if clock == nil {
		clock = pkg.RealClock{}
	}
	return &TokenHandler{
		tokenUsecase:  tokenUsecase,
		clientUsecase: clientUsecase,
		clock:         clock,
	}
}

func (t *TokenHandler) PostToken(c *gin.Context) {
	clientID, clientSecret, err := t.parseBasicAuth(c)
	if err != nil {
		logger.Info(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "client authentication is required"))
		return
	}

	client, err := t.clientUsecase.Authenticate(clientID, clientSecret)
	if err != nil {
		logger.Info(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "invalid client"))
		return
	}

	var request presenter.TokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Info(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusBadRequest, "invalid request"))
		return
	}

	token, err := t.tokenUsecase.Refresh(request.RefreshToken, client.ClientID)
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

func (t *TokenHandler) parseBasicAuth(c *gin.Context) (string, string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", "", errors.New("authorization header is required")
	}

	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Basic") {
		return "", "", errors.New("invalid authorization header")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("invalid basic auth")
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return "", "", errors.New("invalid basic auth")
	}

	return credentials[0], credentials[1], nil
}
