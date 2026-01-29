package integration

import (
	"context"
	"encoding/base64"
	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/api"
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
	apiClient, err := presenter.NewClientWithResponses(baseEndpoint)
	t.Require().NoError(err)

	refreshToken := t.getRefreshToken()
	tokenResponse, err := apiClient.PostTokenWithResponse(context.Background(), presenter.TokenRequest{
		RefreshToken: refreshToken,
	}, t.basicAuthEditor())
	t.Require().NoError(err)
	t.Require().NotNil(tokenResponse.JSON200)
	accessToken := tokenResponse.JSON200.Data.AccessToken

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

func (t *AccountInfoTestSuite) TestPostToken() {
	baseEndpoint := pkg.GetEndpoint("api/v1")
	apiClient, err := presenter.NewClientWithResponses(baseEndpoint)
	t.Require().NoError(err)

	refreshToken := t.getRefreshToken()
	response, err := apiClient.PostTokenWithResponse(context.Background(), presenter.TokenRequest{
		RefreshToken: refreshToken,
	}, t.basicAuthEditor())
	t.Require().NoError(err)
	t.Require().NotNil(response.JSON200)
	t.Assert().Equal(http.StatusOK, response.StatusCode())
	t.Assert().Equal(api.Version, response.JSON200.ApiVersion)
	t.Assert().Equal("Bearer", response.JSON200.Data.TokenType)
	t.Assert().NotEmpty(response.JSON200.Data.AccessToken)
	t.Assert().NotEmpty(response.JSON200.Data.RefreshToken)
	t.Assert().NotEqual("test-refresh-token", response.JSON200.Data.RefreshToken)
	t.Assert().Greater(response.JSON200.Data.ExpiresIn, 0)
	t.Assert().LessOrEqual(response.JSON200.Data.ExpiresIn, 3600)

	var storedToken entity.Token
	err = t.DB.Where("access_token = ?", response.JSON200.Data.AccessToken).Take(&storedToken).Error
	t.Require().NoError(err)
	t.Assert().Equal(response.JSON200.Data.RefreshToken, storedToken.RefreshToken)
	t.Assert().Equal(1, storedToken.CifNo)
	t.Assert().Equal(testClientID, storedToken.ClientID)
	t.Assert().True(storedToken.ExpiresAt.After(time.Now()))
}

func (t *AccountInfoTestSuite) getRefreshToken() string {
	var token entity.Token
	err := t.DB.Select("refresh_token").Where("cif_no = ?", 1).Take(&token).Error
	t.Require().NoError(err)
	return token.RefreshToken
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
	if err := t.DB.Exec("DELETE FROM clients").Error; err != nil {
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

const (
	testClientID     = "client-1"
	testClientSecret = "secret-1"
)

func (t *AccountInfoTestSuite) basicAuthEditor() func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		credentials := base64.StdEncoding.EncodeToString([]byte(testClientID + ":" + testClientSecret))
		req.Header.Set("Authorization", "Basic "+credentials)
		return nil
	}
}

func (t *AccountInfoTestSuite) seedDatabase() error {
	if t.DB == nil {
		return nil
	}

	secretHash, err := pkg.HashString(testClientSecret)
	if err != nil {
		return err
	}

	if err := t.DB.Create(&entity.Client{
		ClientID:     testClientID,
		ClientSecret: secretHash,
		ClientName:   "Test Client",
		Scope:        "read:account_and_transactions",
	}).Error; err != nil {
		return err
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
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		Scopes:       "read:account_and_transactions",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
		CifNo:        1,
		ClientID:     testClientID,
	}).Error; err != nil {
		return err
	}

	return nil
}
