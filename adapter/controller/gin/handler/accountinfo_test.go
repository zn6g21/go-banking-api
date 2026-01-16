package handler

import (
	"encoding/json"
	"errors"
	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/entity"
	"go-banking-api/pkg"
	"go-banking-api/usecase"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAccountInfoUsecase struct {
	mock.Mock
}

func NewMockAccountInfoUsecase() *MockAccountInfoUsecase {
	return &MockAccountInfoUsecase{}
}

func (m *MockAccountInfoUsecase) Get(cifNo int) (*usecase.AccountInfo, error) {
	args := m.Called(cifNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.AccountInfo), args.Error(1)
}

type MockTokenUsecase struct {
	mock.Mock
}

func NewMockTokenUsecase() *MockTokenUsecase {
	return &MockTokenUsecase{}
}

func (m *MockTokenUsecase) Validate(accessTokenFromHeader string, requiredScope string) (*entity.Token, error) {
	args := m.Called(accessTokenFromHeader, requiredScope)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Token), args.Error(1)
}

type AccountInfoHandlerSuite struct {
	suite.Suite
	accountInfoHandler *AccountInfoHandler
}

func TestAccounInfoHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(AccountInfoHandlerSuite))
}

func (suite *AccountInfoHandlerSuite) TestGet() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}

	mockTokenUsecase.On("Validate", "access-token-1", "read:account_and_transactions").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "read:account_and_transactions",
		ExpiresAt:   fixedNow.Add(1 * time.Hour),
		CifNo:       1,
	}, nil)
	mockUsecase.On("Get", 1).Return(&usecase.AccountInfo{
		Status:        entity.AccountStatusActive,
		BranchCode:    "123",
		AccountNumber: "1234567",
		AccountType:   "1",
		Currency:      "JPY",
		Balance:       100000,
		NameKana:      "Tanaka Taro",
		NameKanji:     "田中 太郎",
	}, nil)
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, clock)

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Bearer access-token-1")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var accountResponse presenter.AccountResponse
	err := json.Unmarshal(bodyBytes, &accountResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal(fixedNow, accountResponse.BaseDate.Time)
	suite.Assert().Equal("1234", accountResponse.Data.BankCode)
	suite.Assert().Equal("123", accountResponse.Data.BranchCode)
	suite.Assert().Equal(presenter.Active, accountResponse.Data.Status)
	suite.Assert().Equal("1", accountResponse.Data.AccountType)
	suite.Assert().Equal("1234567", accountResponse.Data.AccountNumber)
	suite.Assert().Equal("JPY", accountResponse.Data.Currency)
	suite.Assert().Equal("100000", accountResponse.Data.Balance)
	suite.Assert().Equal("Tanaka Taro", accountResponse.Data.NameKana)
	suite.Assert().Equal("田中 太郎", accountResponse.Data.NameKanji)

}

func (suite *AccountInfoHandlerSuite) TestGet_MissingAuthorizationHeader() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, pkg.FixedClock{})

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err := json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusUnauthorized, w.Code)
	suite.Assert().Equal(http.StatusUnauthorized, errorResponse.Error.Code)
	suite.Assert().Equal("access token is required", errorResponse.Error.Message)
}

func (suite *AccountInfoHandlerSuite) TestGet_InvalidAuthorizationHeader() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, pkg.FixedClock{})

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Token access-token-1")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err := json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusUnauthorized, w.Code)
	suite.Assert().Equal(http.StatusUnauthorized, errorResponse.Error.Code)
	suite.Assert().Equal("invalid access token", errorResponse.Error.Message)
}

func (suite *AccountInfoHandlerSuite) TestGet_AccountNotFound() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, pkg.FixedClock{})

	mockTokenUsecase.On("Validate", "access-token-1", "read:account_and_transactions").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "read:account_and_transactions",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CifNo:       1,
	}, nil)
	mockUsecase.On("Get", 1).Return(nil, usecase.ErrAccountNotFound)

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Bearer access-token-1")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err := json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusNotFound, w.Code)
	suite.Assert().Equal(http.StatusNotFound, errorResponse.Error.Code)
	suite.Assert().Equal("account not found", errorResponse.Error.Message)
}

func (suite *AccountInfoHandlerSuite) TestGet_AccountInactive() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, pkg.FixedClock{})

	mockTokenUsecase.On("Validate", "access-token-1", "read:account_and_transactions").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "read:account_and_transactions",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CifNo:       1,
	}, nil)
	mockUsecase.On("Get", 1).Return(nil, usecase.ErrAccountInactive)

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Bearer access-token-1")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err := json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusNotFound, w.Code)
	suite.Assert().Equal(http.StatusNotFound, errorResponse.Error.Code)
	suite.Assert().Equal("account not found", errorResponse.Error.Message)
}

func (suite *AccountInfoHandlerSuite) TestGet_TokenValidationError() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, pkg.FixedClock{})

	mockTokenUsecase.On("Validate", "access-token-1", "read:account_and_transactions").Return(nil, errors.New("token invalid"))

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Bearer access-token-1")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err := json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusUnauthorized, w.Code)
	suite.Assert().Equal(http.StatusUnauthorized, errorResponse.Error.Code)
	suite.Assert().Equal("invalid access token", errorResponse.Error.Message)
}

func (suite *AccountInfoHandlerSuite) TestGet_UsecaseError() {
	mockUsecase := NewMockAccountInfoUsecase()
	mockTokenUsecase := NewMockTokenUsecase()
	suite.accountInfoHandler = NewAccountInfoHandler(mockUsecase, mockTokenUsecase, pkg.FixedClock{})

	mockTokenUsecase.On("Validate", "access-token-1", "read:account_and_transactions").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "read:account_and_transactions",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CifNo:       1,
	}, nil)
	mockUsecase.On("Get", 1).Return(&usecase.AccountInfo{}, errors.New("db error"))

	request, _ := http.NewRequest("GET", "/api/v1/accounts", nil)
	request.Header.Set("Authorization", "Bearer access-token-1")
	w := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(w)
	ginContext.Request = request

	suite.accountInfoHandler.GetAccountInformation(ginContext)

	bodyBytes, _ := io.ReadAll(w.Body)
	var errorResponse presenter.ErrorResponse
	err := json.Unmarshal(bodyBytes, &errorResponse)
	suite.Assert().Nil(err)
	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().Equal(http.StatusInternalServerError, errorResponse.Error.Code)
	suite.Assert().Equal("internal server error", errorResponse.Error.Message)
}
