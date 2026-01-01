package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/api"
	"go-banking-api/pkg"
	"go-banking-api/pkg/logger"
	"go-banking-api/usecase"
)

type AccountInfoHandler struct {
	accountInfoUseCase usecase.AccountInfoUsecase
	tokenUsecase       usecase.TokenUsecase
	clock              pkg.Clock
}

func NewAccountInfoHandler(accountInfoUseCase usecase.AccountInfoUsecase, tokenUsecase usecase.TokenUsecase, clock pkg.Clock) *AccountInfoHandler {
	if clock == nil {
		clock = pkg.RealClock{}
	}
	return &AccountInfoHandler{
		accountInfoUseCase: accountInfoUseCase,
		tokenUsecase:       tokenUsecase,
		clock:              clock,
	}
}

func (a *AccountInfoHandler) GetAccountInformation(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		logger.Info("authorization header is required")
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "access token is required"))
		return
	}

	parts := strings.Fields(authorization)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		logger.Info("invalid authorization header")
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "invalid access token"))
		return
	}
	token := parts[1]

	validatedToken, err := a.tokenUsecase.Validate(token, "read:account_and_transactions")
	if err != nil {
		logger.Info(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusUnauthorized, "invalid access token"))
		return
	}

	accountInfo, err := a.accountInfoUseCase.Get(validatedToken.CifNo)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(presenter.NewErrorResponse(http.StatusInternalServerError, "internal server error"))
		return
	}
	c.JSON(http.StatusOK, a.accountInfoToResponse(accountInfo))
}

func (a *AccountInfoHandler) GetTransactionList(c *gin.Context) {
	c.JSON(presenter.NewErrorResponse(http.StatusNotImplemented, "not implemented"))
}

func (a *AccountInfoHandler) accountInfoToResponse(accountInfo *usecase.AccountInfo) *presenter.AccountResponse {
	return &presenter.AccountResponse{
		ApiVersion: api.Version,
		BaseDate:   api.NewBaseDate(a.clock.Now()),
		Data: presenter.Account{
			BankCode:      api.BankCode,
			BranchCode:    accountInfo.BranchCode,
			Status:        presenter.AccountStatus(accountInfo.Status),
			AccountType:   accountInfo.AccountType,
			AccountNumber: accountInfo.AccountNumber,
			Currency:      accountInfo.Currency,
			Balance:       strconv.FormatInt(accountInfo.Balance, 10),
			NameKana:      accountInfo.NameKana,
			NameKanji:     accountInfo.NameKanji,
		},
	}
}
