package handler

import (
	"go-banking-api/entity"

	"github.com/stretchr/testify/mock"
)

type MockTokenUsecase struct {
	mock.Mock
}

func NewMockTokenUsecase() *MockTokenUsecase {
	return &MockTokenUsecase{}
}

func (m *MockTokenUsecase) Validate(accessTokenFromHeader string, requiredScope string) (*entity.Token, error) {
	args := m.Called(accessTokenFromHeader, requiredScope)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Token), args.Error(1)
}

func (m *MockTokenUsecase) Refresh(refreshToken string, clientID string) (*entity.Token, error) {
	args := m.Called(refreshToken, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Token), args.Error(1)
}

type MockClientUsecase struct {
	mock.Mock
}

func NewMockClientUsecase() *MockClientUsecase {
	return &MockClientUsecase{}
}

func (m *MockClientUsecase) Authenticate(clientID string, clientSecret string) (*entity.Client, error) {
	args := m.Called(clientID, clientSecret)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Client), args.Error(1)
}
