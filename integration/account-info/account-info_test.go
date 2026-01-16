package integration

import (
	"context"
	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/entity"
	"go-banking-api/infrastructure/database"
	"go-banking-api/pkg"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AccountInfoTestSuite struct {
	suite.Suite
	DB *gorm.DB
}

func TestAccountInfoTestSuite(t *testing.T) {
	suite.Run(t, new(AccountInfoTestSuite))
}

func (t *AccountInfoTestSuite) SetupSuite() {
	if os.Getenv("APP_ENV") != "integration" {
		t.T().Skip("integration test: set APP_ENV=integration and run via make integration-test")
	}

	db, err := database.NewDatabaseSQLFactory(database.InstanceMySQL)
	t.Require().NoError(err)
	t.DB = db

	t.Require().NoError(t.waitForHealth(10 * time.Second))

	t.Require().NoError(t.cleanupDatabase())
	t.Require().NoError(t.seedDatabase())
}

func (t *AccountInfoTestSuite) TearDownSuite() {
	_ = t.cleanupDatabase()
	if t.DB != nil {
		sqlDB, err := t.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}
}

func (t *AccountInfoTestSuite) TestGetAccountInfo() {
	baseEndpoint := pkg.GetEndpoint("api/v1")
	apiClient, _ := presenter.NewClientWithResponses(baseEndpoint)
	const accessToken = "test-access-token"

	authEditor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}

	getResponse, err := apiClient.GetAccountInformationWithResponse(context.Background(), authEditor)
	t.Assert().Nil(err)
	t.Assert().Equal(http.StatusOK, getResponse.StatusCode())
	t.Assert().Equal("1234", getResponse.JSON200.Data.BankCode)
	t.Assert().Equal("123", getResponse.JSON200.Data.BranchCode)
	t.Assert().Equal(presenter.Active, getResponse.JSON200.Data.Status)
	t.Assert().Equal("1", getResponse.JSON200.Data.AccountType)
	t.Assert().Equal("1234567", getResponse.JSON200.Data.AccountNumber)
	t.Assert().Equal("JPY", getResponse.JSON200.Data.Currency)
	t.Assert().Equal("100000", getResponse.JSON200.Data.Balance)
	t.Assert().Equal("Tanaka Taro", getResponse.JSON200.Data.NameKana)
	t.Assert().Equal("田中 太郎", getResponse.JSON200.Data.NameKanji)

}

func (t *AccountInfoTestSuite) waitForHealth(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	url := pkg.GetEndpoint("health")
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return context.DeadlineExceeded
}

func (t *AccountInfoTestSuite) cleanupDatabase() error {
	if t.DB == nil {
		return nil
	}
	if err := t.DB.Exec("DELETE FROM tokens").Error; err != nil {
		return err
	}
	if err := t.DB.Exec("DELETE FROM accounts").Error; err != nil {
		return err
	}
	if err := t.DB.Exec("DELETE FROM customers").Error; err != nil {
		return err
	}
	return nil
}

func (t *AccountInfoTestSuite) seedDatabase() error {
	if t.DB == nil {
		return nil
	}

	if err := t.DB.Create(&entity.Customer{
		CifNo:      1,
		NameKana:   "Tanaka Taro",
		NameKanji:  "田中 太郎",
		BirthDate:  pkg.Str2time("1990-01-01"),
		Prefecture: "Tokyo",
		City:       "Chiyoda",
		Town:       "Kanda",
		Street:     "1-1-1",
		Building:   "",
		Room:       "",
		Email:      "tanaka@example.com",
		Phone:      "000-0000-0000",
	}).Error; err != nil {
		return err
	}

	if err := t.DB.Create(&entity.Account{
		Id:            1,
		CifNo:         1,
		Status:        entity.AccountStatusActive,
		BranchCode:    "123",
		AccountNumber: "1234567",
		AccountType:   "1",
		Currency:      "JPY",
		Balance:       int64(100000),
	}).Error; err != nil {
		return err
	}

	if err := t.DB.Create(&entity.Token{
		AccessToken: "test-access-token",
		Scopes:      "read:account_and_transactions",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		CifNo:       1,
	}).Error; err != nil {
		return err
	}

	return nil
}
