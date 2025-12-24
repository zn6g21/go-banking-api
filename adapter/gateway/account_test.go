package gateway_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"

	"go-banking-api/adapter/gateway"
	"go-banking-api/entity"
	"go-banking-api/pkg"
	"go-banking-api/pkg/tester"
)

type AccountRepositoryTestSuite struct {
	tester.DBSQLiteSuite
	repository gateway.AccountRepository
}

func TestAccountRepositorySuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}

func (suite *AccountRepositoryTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewAccountRepository(suite.DBSQLiteSuite.DB)
}

func (suite *AccountRepositoryTestSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewAccountRepository(mockGormDB)
	return mock
}

func (suite *AccountRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewAccountRepository(suite.DB)
}

func (suite *AccountRepositoryTestSuite) TestAccountRepositoryGet() {
	now := pkg.Str2time("2025-12-02")
	paramAccount := entity.Account{
		Id:            1,
		CifNo:         1,
		Status:        entity.AccountStatusActive,
		BranchCode:    "001",
		AccountNumber: "1234567",
		AccountType:   "1",
		Currency:      "JPY",
		Balance:       int64(10000),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	suite.DB.Create(&paramAccount)
	got, err := suite.repository.Get(paramAccount.CifNo)
	suite.Assert().Nil(err)
	suite.Assert().Equal(paramAccount, *got)
}

func (suite *AccountRepositoryTestSuite) TestAccountGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE cif_no = ? LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	account, err := suite.repository.Get(1)
	suite.Assert().Nil(account)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}
