package usecase

import (
	"errors"

	"go-banking-api/adapter/gateway"
	"go-banking-api/entity"
	"go-banking-api/pkg"

	"gorm.io/gorm"
)

var (
	ErrClientIDRequired     = errors.New("client id is required")
	ErrClientSecretRequired = errors.New("client secret is required")
	ErrInvalidClient        = errors.New("invalid client")
)

type ClientUsecase interface {
	Authenticate(clientID string, clientSecret string) (*entity.Client, error)
}

type clientUsecase struct {
	clientRepository gateway.ClientRepository
}

func NewClientUsecase(clientRepository gateway.ClientRepository) *clientUsecase {
	return &clientUsecase{clientRepository: clientRepository}
}

func (c *clientUsecase) Authenticate(clientID string, clientSecret string) (*entity.Client, error) {
	if clientID == "" {
		return nil, ErrClientIDRequired
	}
	if clientSecret == "" {
		return nil, ErrClientSecretRequired
	}

	client, err := c.clientRepository.Get(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidClient
		}
		return nil, err
	}

	if !pkg.CompareHash(client.ClientSecret, clientSecret) {
		return nil, ErrInvalidClient
	}

	return client, nil
}
