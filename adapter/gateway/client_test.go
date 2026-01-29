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

type ClientRepositoryTestSuite struct {
	tester.DBSQLiteSuite
	repository gateway.ClientRepository
}

func TestClientRepositorySuite(t *testing.T) {
	suite.Run(t, new(ClientRepositoryTestSuite))
}

func (suite *ClientRepositoryTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewClientRepository(suite.DB)
}

func (suite *ClientRepositoryTestSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewClientRepository(mockGormDB)
	return mock
}

func (suite *ClientRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewClientRepository(suite.DB)
}

func (suite *ClientRepositoryTestSuite) TestClientRepositoryGet() {
	secretHash, err := pkg.HashString("secret-1")
	suite.Require().NoError(err)

	paramClient := entity.Client{
		ClientID:     "client-1",
		ClientSecret: secretHash,
		ClientName:   "Test Client",
		Scope:        "read:account_and_transactions",
	}

	suite.DB.Create(&paramClient)
	got, err := suite.repository.Get(paramClient.ClientID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(paramClient, *got)
}

func (suite *ClientRepositoryTestSuite) TestClientGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `clients` WHERE client_id = ? LIMIT ?")).
		WithArgs("client-1", 1).
		WillReturnError(errors.New("get error"))

	client, err := suite.repository.Get("client-1")
	suite.Assert().Nil(client)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}
