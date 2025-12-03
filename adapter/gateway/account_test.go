package gateway_test

import (
	"testing"

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

func (suite *AccountRepositoryTestSuite) TestAccount() {
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
	got, err := suite.repository.Get(paramAccount.Id)
	suite.Assert().Nil(err)
	suite.Assert().Equal(paramAccount, *got)
}
