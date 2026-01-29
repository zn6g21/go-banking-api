package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"go-banking-api/entity"
	"go-banking-api/pkg"
)

type mockClientRepository struct {
	mock.Mock
}

func NewMockClientRepository() *mockClientRepository {
	return &mockClientRepository{}
}

func (m *mockClientRepository) Get(clientID string) (*entity.Client, error) {
	args := m.Called(clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Client), args.Error(1)
}

type ClientUsecaseSuite struct {
	suite.Suite
	clientUsecase *clientUsecase
}

func TestClientUsecaseSuite(t *testing.T) {
	suite.Run(t, new(ClientUsecaseSuite))
}

func (suite *ClientUsecaseSuite) TestAuthenticateSuccess() {
	mockClientRepository := NewMockClientRepository()
	suite.clientUsecase = NewClientUsecase(mockClientRepository)

	secretHash, err := pkg.HashString("secret-1")
	suite.Require().NoError(err)

	mockClientRepository.On("Get", "client-1").Return(&entity.Client{
		ClientID:     "client-1",
		ClientSecret: secretHash,
		ClientName:   "Test Client",
		Scope:        "read:account_and_transactions",
	}, nil)

	client, err := suite.clientUsecase.Authenticate("client-1", "secret-1")
	suite.Assert().Nil(err)
	suite.Assert().Equal("client-1", client.ClientID)
}

func (suite *ClientUsecaseSuite) TestAuthenticateMissingClientID() {
	mockClientRepository := NewMockClientRepository()
	suite.clientUsecase = NewClientUsecase(mockClientRepository)

	client, err := suite.clientUsecase.Authenticate("", "secret-1")
	suite.Assert().Nil(client)
	suite.Assert().ErrorIs(err, ErrClientIDRequired)
}

func (suite *ClientUsecaseSuite) TestAuthenticateMissingClientSecret() {
	mockClientRepository := NewMockClientRepository()
	suite.clientUsecase = NewClientUsecase(mockClientRepository)

	client, err := suite.clientUsecase.Authenticate("client-1", "")
	suite.Assert().Nil(client)
	suite.Assert().ErrorIs(err, ErrClientSecretRequired)
}

func (suite *ClientUsecaseSuite) TestAuthenticateNotFound() {
	mockClientRepository := NewMockClientRepository()
	suite.clientUsecase = NewClientUsecase(mockClientRepository)

	mockClientRepository.On("Get", "client-1").Return(nil, gorm.ErrRecordNotFound)

	client, err := suite.clientUsecase.Authenticate("client-1", "secret-1")
	suite.Assert().Nil(client)
	suite.Assert().ErrorIs(err, ErrInvalidClient)
}

func (suite *ClientUsecaseSuite) TestAuthenticateInvalidSecret() {
	mockClientRepository := NewMockClientRepository()
	suite.clientUsecase = NewClientUsecase(mockClientRepository)

	secretHash, err := pkg.HashString("secret-1")
	suite.Require().NoError(err)

	mockClientRepository.On("Get", "client-1").Return(&entity.Client{
		ClientID:     "client-1",
		ClientSecret: secretHash,
		ClientName:   "Test Client",
		Scope:        "read:account_and_transactions",
	}, nil)

	client, err := suite.clientUsecase.Authenticate("client-1", "secret-2")
	suite.Assert().Nil(client)
	suite.Assert().ErrorIs(err, ErrInvalidClient)
}

func (suite *ClientUsecaseSuite) TestAuthenticateRepositoryError() {
	mockClientRepository := NewMockClientRepository()
	suite.clientUsecase = NewClientUsecase(mockClientRepository)

	mockClientRepository.On("Get", "client-1").Return(nil, errors.New("db error"))

	client, err := suite.clientUsecase.Authenticate("client-1", "secret-1")
	suite.Assert().Nil(client)
	suite.Assert().Equal("db error", err.Error())
}
