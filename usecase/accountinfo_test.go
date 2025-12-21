package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go-banking-api/entity"
)

type mockCustomerRepository struct {
	mock.Mock
}

type mockAccountRepository struct {
	mock.Mock
}

func NewMockCustomerRepository() *mockCustomerRepository {
	return &mockCustomerRepository{}
}

func NewMockAccountRepository() *mockAccountRepository {
	return &mockAccountRepository{}
}

func (m *mockCustomerRepository) Get(cifNo int) (*entity.Customer, error) {
	args := m.Called(cifNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Customer), args.Error(1)
}

func (m *mockAccountRepository) Get(cifNo int) (*entity.Account, error) {
	args := m.Called(cifNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}

type AccountInfoUseCaseSuite struct {
	suite.Suite
	accountInfoUseCase *accountInfoUsecase
}

func TestAccountInfoUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AccountInfoUseCaseSuite))
}

func (suite *AccountInfoUseCaseSuite) SetupSuite() {
}

func (suite *AccountInfoUseCaseSuite) TestGet() {
	nameKana := "Taro Tanaka"
	nameKanji := "田中 太郎"
	status := entity.AccountStatusActive
	branchCode := "123"
	accountNumber := "1234567"
	accountType := "1"
	currency := "JPY"
	balance := int64(10000)
	mockCustomerRepository := NewMockCustomerRepository()
	mockAccountRepository := NewMockAccountRepository()
	suite.accountInfoUseCase = NewAccountInfoUsecase(mockCustomerRepository, mockAccountRepository)
	mockCustomerRepository.On("Get", 1).Return(&entity.Customer{
		NameKana:  nameKana,
		NameKanji: nameKanji,
	}, nil)
	mockAccountRepository.On("Get", 1).Return(&entity.Account{
		Status:        status,
		BranchCode:    branchCode,
		AccountNumber: accountNumber,
		AccountType:   accountType,
		Currency:      currency,
		Balance:       balance,
	}, nil)

	accountInfo, err := suite.accountInfoUseCase.Get(1)
	suite.Assert().Nil(err)
	suite.Assert().Equal(&AccountInfo{
		NameKana:      nameKana,
		NameKanji:     nameKanji,
		Status:        status,
		BranchCode:    branchCode,
		AccountNumber: accountNumber,
		AccountType:   accountType,
		Currency:      currency,
		Balance:       balance,
	}, accountInfo)

}

func (suite *AccountInfoUseCaseSuite) TestGetCustomerRepositoryError() {
	expectedErr := errors.New("customer error")
	mockCustomerRepository := NewMockCustomerRepository()
	mockAccountRepository := NewMockAccountRepository()
	suite.accountInfoUseCase = NewAccountInfoUsecase(mockCustomerRepository, mockAccountRepository)

	mockCustomerRepository.On("Get", 1).Return(nil, expectedErr)

	accountInfo, err := suite.accountInfoUseCase.Get(1)
	suite.Assert().Nil(accountInfo)
	suite.Assert().Equal(expectedErr, err)
}

func (suite *AccountInfoUseCaseSuite) TestGetAccountRepositoryError() {
	expectedErr := errors.New("account error")
	mockCustomerRepository := NewMockCustomerRepository()
	mockAccountRepository := NewMockAccountRepository()
	suite.accountInfoUseCase = NewAccountInfoUsecase(mockCustomerRepository, mockAccountRepository)

	mockCustomerRepository.On("Get", 1).Return(&entity.Customer{
		NameKana:  "Taro Tanaka",
		NameKanji: "田中 太郎",
	}, nil)
	mockAccountRepository.On("Get", 1).Return(nil, expectedErr)

	accountInfo, err := suite.accountInfoUseCase.Get(1)
	suite.Assert().Nil(accountInfo)
	suite.Assert().Equal(expectedErr, err)
}

func (suite *AccountInfoUseCaseSuite) TestGetAccountNotActive() {
	expectedErr := errors.New("account is not active")
	mockCustomerRepository := NewMockCustomerRepository()
	mockAccountRepository := NewMockAccountRepository()
	suite.accountInfoUseCase = NewAccountInfoUsecase(mockCustomerRepository, mockAccountRepository)

	mockCustomerRepository.On("Get", 1).Return(&entity.Customer{
		NameKana:  "Taro Tanaka",
		NameKanji: "田中 太郎",
	}, nil)
	mockAccountRepository.On("Get", 1).Return(&entity.Account{
		Status: entity.AccountStatusClosed,
	}, nil)

	accountInfo, err := suite.accountInfoUseCase.Get(1)
	suite.Assert().Nil(accountInfo)
	suite.Assert().Equal(expectedErr, err)
}
