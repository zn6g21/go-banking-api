package usecase

import (
	"errors"

	"go-banking-api/adapter/gateway"
	"go-banking-api/entity"
	"go-banking-api/pkg"

	"gorm.io/gorm"
)

type TokenUsecase interface {
	Validate(accessTokenFromHeader string, requiredScope string) (*entity.Token, error)
}

type tokenUsecase struct {
	tokenRepository gateway.TokenRepository
	clock           pkg.Clock
}

func NewTokenUsecase(tokenRepository gateway.TokenRepository, clock pkg.Clock) *tokenUsecase {
	if clock == nil {
		clock = pkg.RealClock{}
	}
	return &tokenUsecase{tokenRepository: tokenRepository, clock: clock}
}

func (t *tokenUsecase) Validate(accessTokenFromHeader string, requiredScope string) (*entity.Token, error) {
	if accessTokenFromHeader == "" {
		return nil, errors.New("access token is required")
	}

	storedToken, err := t.tokenRepository.Get(accessTokenFromHeader)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid access token")
		}
		return nil, err
	}

	if storedToken.IsExpired(t.clock) {
		return nil, errors.New("token expired")
	}

	if !storedToken.HasScope(requiredScope) {
		return nil, errors.New("invalid scope")
	}

	return storedToken, nil
}
