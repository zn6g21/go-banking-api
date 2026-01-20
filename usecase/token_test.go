package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-banking-api/entity"
	"go-banking-api/pkg"
)

type mockTokenRepository struct {
	mock.Mock
}

func NewMockTokenRepository() *mockTokenRepository {
	return &mockTokenRepository{}
}

func (m *mockTokenRepository) Get(token string) (*entity.Token, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Token), args.Error(1)
}

func (m *mockTokenRepository) UpdateByRefreshToken(refreshToken string, accessToken string, newRefreshToken string, expiresAt time.Time) error {
	args := m.Called(refreshToken, accessToken, newRefreshToken, expiresAt)
	return args.Error(0)
}

type TokenUsecaseSuite struct {
	suite.Suite
	tokenUsecase *tokenUsecase
}

func TestTokenUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TokenUsecaseSuite))
}

func (suite *TokenUsecaseSuite) TestValidate() {
	mockTokenRepository := NewMockTokenRepository()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, clock)

	expiresAt := fixedNow.Add(1 * time.Hour)
	requiredScope := "read:account_and_transactions"
	mockTokenRepository.On("Get", "access-token-1").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "read:account_and_transactions write:transfer",
		ExpiresAt:   expiresAt,
		CifNo:       1,
	}, nil)

	token, err := suite.tokenUsecase.Validate("access-token-1", requiredScope)
	suite.Assert().Nil(err)
	suite.Assert().Equal("access-token-1", token.AccessToken)
}

func (suite *TokenUsecaseSuite) TestValidateEmptyAccessToken() {
	mockTokenRepository := NewMockTokenRepository()
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, pkg.FixedClock{T: time.Now()})

	token, err := suite.tokenUsecase.Validate("", "read:account_and_transactions")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("access token is required", err.Error())
}

func (suite *TokenUsecaseSuite) TestValidateInvalidAccessToken() {
	mockTokenRepository := NewMockTokenRepository()
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, pkg.FixedClock{T: time.Now()})

	mockTokenRepository.On("Get", "access-token-1").Return(nil, gorm.ErrRecordNotFound)

	token, err := suite.tokenUsecase.Validate("access-token-1", "read:account_and_transactions")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("invalid access token", err.Error())
}

func (suite *TokenUsecaseSuite) TestValidateRepositoryError() {
	mockTokenRepository := NewMockTokenRepository()
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, pkg.FixedClock{T: time.Now()})

	mockTokenRepository.On("Get", "access-token-1").Return(nil, errors.New("get error"))

	token, err := suite.tokenUsecase.Validate("access-token-1", "read:account_and_transactions")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *TokenUsecaseSuite) TestValidateExpired() {
	mockTokenRepository := NewMockTokenRepository()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, clock)

	mockTokenRepository.On("Get", "access-token-1").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "read:account_and_transactions",
		ExpiresAt:   fixedNow.Add(-1 * time.Hour),
		CifNo:       1,
	}, nil)

	token, err := suite.tokenUsecase.Validate("access-token-1", "read:account_and_transactions")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("token expired", err.Error())
}

func (suite *TokenUsecaseSuite) TestValidateInvalidScope() {
	mockTokenRepository := NewMockTokenRepository()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, clock)

	mockTokenRepository.On("Get", "access-token-1").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "write:transfer",
		ExpiresAt:   fixedNow.Add(1 * time.Hour),
		CifNo:       1,
	}, nil)

	token, err := suite.tokenUsecase.Validate("access-token-1", "read:account_and_transactions")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("invalid scope", err.Error())
}

func (suite *TokenUsecaseSuite) TestValidateScopeNotRequired() {
	mockTokenRepository := NewMockTokenRepository()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, clock)

	mockTokenRepository.On("Get", "access-token-1").Return(&entity.Token{
		AccessToken: "access-token-1",
		Scopes:      "",
		ExpiresAt:   fixedNow.Add(1 * time.Hour),
		CifNo:       1,
	}, nil)

	token, err := suite.tokenUsecase.Validate("access-token-1", "")
	suite.Assert().Nil(err)
	suite.Assert().Equal("access-token-1", token.AccessToken)
}

func (suite *TokenUsecaseSuite) TestRefresh() {
	mockTokenRepository := NewMockTokenRepository()
	fixedNow := time.Date(2025, 12, 21, 0, 0, 0, 0, time.UTC)
	clock := pkg.FixedClock{T: fixedNow}
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, clock)

	expectedExpiresAt := fixedNow.Add(1 * time.Hour)
	mockTokenRepository.On(
		"UpdateByRefreshToken",
		"refresh-token-1",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		expectedExpiresAt,
	).Return(nil)

	token, err := suite.tokenUsecase.Refresh("refresh-token-1")
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(token.AccessToken)
	suite.Assert().NotEmpty(token.RefreshToken)
	suite.Assert().Equal(expectedExpiresAt, token.ExpiresAt)
}

func (suite *TokenUsecaseSuite) TestRefreshEmptyRefreshToken() {
	mockTokenRepository := NewMockTokenRepository()
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, pkg.FixedClock{T: time.Now()})

	token, err := suite.tokenUsecase.Refresh("")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("refresh token is required", err.Error())
}

func (suite *TokenUsecaseSuite) TestRefreshInvalidRefreshToken() {
	mockTokenRepository := NewMockTokenRepository()
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, pkg.FixedClock{T: time.Now()})

	mockTokenRepository.On("UpdateByRefreshToken", "refresh-token-1", mock.Anything, mock.Anything, mock.Anything).
		Return(gorm.ErrRecordNotFound)

	token, err := suite.tokenUsecase.Refresh("refresh-token-1")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("invalid refresh token", err.Error())
}

func (suite *TokenUsecaseSuite) TestRefreshUpdateError() {
	mockTokenRepository := NewMockTokenRepository()
	suite.tokenUsecase = NewTokenUsecase(mockTokenRepository, pkg.FixedClock{T: time.Now()})

	mockTokenRepository.On("UpdateByRefreshToken", "refresh-token-1", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("update error"))

	token, err := suite.tokenUsecase.Refresh("refresh-token-1")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}
