package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/api"
	"go-banking-api/entity"
	"go-banking-api/pkg"
	"go-banking-api/usecase"
)

type TokenHandlerSuite struct {
	suite.Suite
	tokenHandler *TokenHandler
}

func TestTokenHandlerSuite(t *testing.T) {
	suite.Run(t, new(TokenHandlerSuite))
}

func (suite *TokenHandlerSuite) TestPostTokenSuccess() {
	mockTokenUsecase := NewMockTokenUsecase()
	mockClientUsecase := NewMockClientUsecase()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}
	expectedToken := &entity.Token{
		AccessToken:  "access-token-1",
		RefreshToken: "refresh-token-2",
		ExpiresAt:    fixedNow.Add(1 * time.Hour),
	}
	mockClientUsecase.On("Authenticate", "client-1", "secret-1").Return(&entity.Client{ClientID: "client-1"}, nil)
	mockTokenUsecase.On("Refresh", "refresh-token-1", "client-1").Return(expectedToken, nil)

	suite.tokenHandler = NewTokenHandler(mockTokenUsecase, mockClientUsecase, clock)

	body, err := json.Marshal(presenter.TokenRequest{RefreshToken: "refresh-token-1"})
	request, err := http.NewRequest("POST", "/api/v1/token", bytes.NewReader(body))
	request.SetBasicAuth("client-1", "secret-1")
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.tokenHandler.PostToken(ginContext)

	bodyBytes, err := io.ReadAll(w.Body)
	var tokenResponse presenter.TokenResponse
	err = json.Unmarshal(bodyBytes, &tokenResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal(api.Version, tokenResponse.ApiVersion)
	suite.Assert().Equal(expectedToken.AccessToken, tokenResponse.Data.AccessToken)
	suite.Assert().Equal(expectedToken.RefreshToken, tokenResponse.Data.RefreshToken)
	suite.Assert().Equal("Bearer", tokenResponse.Data.TokenType)
	suite.Assert().Equal(3600, tokenResponse.Data.ExpiresIn)
}

func (suite *TokenHandlerSuite) TestPostTokenMissingRefreshToken() {
	mockTokenUsecase := NewMockTokenUsecase()
	mockClientUsecase := NewMockClientUsecase()
	mockClientUsecase.On("Authenticate", "client-1", "secret-1").Return(&entity.Client{ClientID: "client-1"}, nil)
	mockTokenUsecase.On("Refresh", "", "client-1").Return(nil, usecase.ErrRefreshTokenRequired)

	suite.tokenHandler = NewTokenHandler(mockTokenUsecase, mockClientUsecase, pkg.FixedClock{T: time.Now()})

	request, err := http.NewRequest("POST", "/api/v1/token", bytes.NewReader([]byte(`{}`)))
	request.SetBasicAuth("client-1", "secret-1")
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.tokenHandler.PostToken(ginContext)

	bodyBytes, err := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err = json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().Equal(http.StatusBadRequest, errorResponse.Error.Code)
	suite.Assert().Equal("refresh token is required", errorResponse.Error.Message)
}

func (suite *TokenHandlerSuite) TestPostTokenInvalidRefreshToken() {
	mockTokenUsecase := NewMockTokenUsecase()
	mockClientUsecase := NewMockClientUsecase()
	mockClientUsecase.On("Authenticate", "client-1", "secret-1").Return(&entity.Client{ClientID: "client-1"}, nil)
	mockTokenUsecase.On("Refresh", "refresh-token-1", "client-1").Return(nil, usecase.ErrInvalidRefreshToken)
	suite.tokenHandler = NewTokenHandler(mockTokenUsecase, mockClientUsecase, pkg.FixedClock{T: time.Now()})

	body, err := json.Marshal(presenter.TokenRequest{RefreshToken: "refresh-token-1"})
	request, err := http.NewRequest("POST", "/api/v1/token", bytes.NewReader(body))
	request.SetBasicAuth("client-1", "secret-1")
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.tokenHandler.PostToken(ginContext)

	bodyBytes, err := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err = json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusUnauthorized, w.Code)
	suite.Assert().Equal(http.StatusUnauthorized, errorResponse.Error.Code)
	suite.Assert().Equal("invalid refresh token", errorResponse.Error.Message)
}

func (suite *TokenHandlerSuite) TestPostTokenUsecaseError() {
	mockTokenUsecase := NewMockTokenUsecase()
	mockClientUsecase := NewMockClientUsecase()
	mockClientUsecase.On("Authenticate", "client-1", "secret-1").Return(&entity.Client{ClientID: "client-1"}, nil)
	mockTokenUsecase.On("Refresh", "refresh-token-1", "client-1").Return(nil, errors.New("db error"))
	suite.tokenHandler = NewTokenHandler(mockTokenUsecase, mockClientUsecase, pkg.FixedClock{T: time.Now()})

	body, err := json.Marshal(presenter.TokenRequest{RefreshToken: "refresh-token-1"})
	request, err := http.NewRequest("POST", "/api/v1/token", bytes.NewReader(body))
	request.SetBasicAuth("client-1", "secret-1")
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.tokenHandler.PostToken(ginContext)

	bodyBytes, err := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err = json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().Equal(http.StatusInternalServerError, errorResponse.Error.Code)
	suite.Assert().Equal("internal server error", errorResponse.Error.Message)
}
