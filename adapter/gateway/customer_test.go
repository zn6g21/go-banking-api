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

type CustomerRepositoryTestSuite struct {
	tester.DBSQLiteSuite
	repository gateway.CustomerRepository
}

func TestCustomerRepositorySuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}

func (suite *CustomerRepositoryTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.repository = gateway.NewCustomerRepository(suite.DBSQLiteSuite.DB)
}

func (suite *CustomerRepositoryTestSuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := tester.MockDB()
	suite.repository = gateway.NewCustomerRepository(mockGormDB)
	return mock
}

func (suite *CustomerRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.repository = gateway.NewCustomerRepository(suite.DB)
}

func (suite *CustomerRepositoryTestSuite) TestCustomer() {
	now := pkg.Str2time("2025-12-02")
	birthDate := pkg.Str2time("1990-01-01")
	paramCustomer := entity.Customer{
		CifNo:      1,
		NameKana:   "Taro Tanaka",
		NameKanji:  "田中 太郎",
		BirthDate:  birthDate,
		Prefecture: "Tokyo",
		City:       "Bunkyo",
		Town:       "Kouraku",
		Street:     "1-1-1",
		Building:   "CivicCenter",
		Room:       "101",
		Email:      "tarou.tanaka@example.com",
		Phone:      "09012345678",
		CreatedAt:  now,
	}

	suite.DB.Create(&paramCustomer)
	got, err := suite.repository.Get(paramCustomer.CifNo)
	suite.Assert().Nil(err)
	suite.Assert().Equal(paramCustomer, *got)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customers` WHERE `customers`.`cif_no` = ? LIMIT ?")).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	customer, err := suite.repository.Get(1)
	suite.Assert().Nil(customer)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}
