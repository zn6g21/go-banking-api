package gateway_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-banking-api/adapter/gateway"
	"go-banking-api/entity"
	"go-banking-api/pkg"
	"go-banking-api/pkg/tester"
)

type TokenRepositoryTestSuite struct {
	tester.DBSQLiteSuite
	repository gateway.TokenRepository
}

func TestTokenRepositorySuite(t *testing.T) {
	suite.Run(t, new(TokenRepositoryTestSuite))
}

func (suite *TokenRepositoryTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewTokenRepository(suite.DB)
}

func (suite *TokenRepositoryTestSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewTokenRepository(mockGormDB)
	return mock
}

func (suite *TokenRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewTokenRepository(suite.DB)
}

func (suite *TokenRepositoryTestSuite) TestTokenRepositoryGet() {
	expiresAt := pkg.Str2time("2025-12-02")
	paramToken := entity.Token{
		AccessToken:  "access-token-1",
		RefreshToken: "refresh-token-1",
		Scopes:       "read:account_and_transactions",
		ExpiresAt:    expiresAt,
		CifNo:        1,
	}

	suite.DB.Create(&paramToken)
	got, err := suite.repository.Get(paramToken.AccessToken)
	suite.Assert().Nil(err)
	suite.Assert().Equal(paramToken, *got)
}

func (suite *TokenRepositoryTestSuite) TestTokenRepositoryUpdateByRefreshToken() {
	expiresAt := pkg.Str2time("2025-12-02")
	paramToken := entity.Token{
		AccessToken:  "access-token-1",
		RefreshToken: "refresh-token-1",
		Scopes:       "read:account_and_transactions",
		ExpiresAt:    expiresAt,
		CifNo:        1,
	}

	suite.DB.Create(&paramToken)

	newExpiresAt := pkg.Str2time("2026-01-01")
	err := suite.repository.UpdateByRefreshToken("refresh-token-1", "access-token-2", "refresh-token-2", newExpiresAt)
	suite.Assert().Nil(err)

	got, err := suite.repository.Get("access-token-2")
	suite.Assert().Nil(err)
	suite.Assert().Equal("refresh-token-2", got.RefreshToken)
	suite.Assert().Equal(newExpiresAt, got.ExpiresAt)
	suite.Assert().Equal(paramToken.Scopes, got.Scopes)
	suite.Assert().Equal(paramToken.CifNo, got.CifNo)
}

func (suite *TokenRepositoryTestSuite) TestTokenRepositoryUpdateByRefreshTokenNotFound() {
	err := suite.repository.UpdateByRefreshToken("missing-refresh-token", "access-token-2", "refresh-token-2", time.Now())
	suite.Assert().NotNil(err)
	suite.Assert().True(errors.Is(err, gorm.ErrRecordNotFound))
}

func (suite *TokenRepositoryTestSuite) TestTokenRepositoryUpdateByRefreshTokenError() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE `tokens` SET")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("update error"))
	mockDB.ExpectRollback()

	err := suite.repository.UpdateByRefreshToken("refresh-token-1", "access-token-2", "refresh-token-2", time.Now())
	suite.Assert().NotNil(err)
	suite.Assert().Equal("update error", err.Error())
}

func (suite *TokenRepositoryTestSuite) TestTokenGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tokens` WHERE access_token = ? LIMIT ?")).WithArgs("access-token-1", 1).WillReturnError(errors.New("get error"))

	token, err := suite.repository.Get("access-token-1")
	suite.Assert().Nil(token)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}
