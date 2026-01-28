package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"go-banking-api/adapter/gateway"
	"go-banking-api/entity"
	"go-banking-api/pkg"

	"gorm.io/gorm"
)

type TokenUsecase interface {
	Validate(accessTokenFromHeader string, requiredScope string) (*entity.Token, error)
	Refresh(refreshToken string, clientID string) (*entity.Token, error)
}

type tokenUsecase struct {
	tokenRepository gateway.TokenRepository
	clock           pkg.Clock
}

const accessTokenTTL = time.Hour

var (
	ErrRefreshTokenRequired = errors.New("refresh token is required")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

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

func (t *tokenUsecase) Refresh(refreshToken string, clientID string) (*entity.Token, error) {
	if refreshToken == "" {
		return nil, ErrRefreshTokenRequired
	}
	if clientID == "" {
		return nil, ErrInvalidRefreshToken
	}

	storedToken, err := t.tokenRepository.GetByRefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}
	if storedToken.ClientID != clientID {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, err := generateToken()
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := t.clock.Now().Add(accessTokenTTL)
	if err := t.tokenRepository.UpdateByRefreshToken(refreshToken, accessToken, newRefreshToken, expiresAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}

	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ClientID:     clientID,
		ExpiresAt:    expiresAt,
	}, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
