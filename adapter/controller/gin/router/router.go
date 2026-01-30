package router

import (
	"encoding/json"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"go-banking-api/adapter/controller/gin/handler"
	"go-banking-api/adapter/controller/gin/middleware"
	"go-banking-api/adapter/controller/gin/presenter"
	"go-banking-api/adapter/gateway"
	"go-banking-api/pkg"
	"go-banking-api/pkg/logger"
	"go-banking-api/usecase"
)

func setupSwagger(router *gin.Engine) (*openapi3.T, error) {
	swagger, err := presenter.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := pkg.GetEnvDefault("APP_ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InfoInstanceName, SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	swagger, err := setupSwagger(router)
	if err != nil {
		logger.Warn(err.Error())
		return nil, err
	}

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	router.GET("/health", handler.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			v1.Use(ginMiddleware.OapiRequestValidatorWithOptions(
				swagger,
				&ginMiddleware.Options{
					Options: openapi3filter.Options{
						AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
					},
				},
			))

			customerRepository := gateway.NewCustomerRepository(db)
			accountRepository := gateway.NewAccountRepository(db)
			clientRepository := gateway.NewClientRepository(db)
			tokenRepository := gateway.NewTokenRepository(db)
			clock := pkg.RealClock{}
			tokenUsecase := usecase.NewTokenUsecase(tokenRepository, clock)
			clientUsecase := usecase.NewClientUsecase(clientRepository)
			accountInfoUseCase := usecase.NewAccountInfoUsecase(customerRepository, accountRepository)
			accountInfoHandler := handler.NewAccountInfoHandler(accountInfoUseCase, tokenUsecase, clock)
			tokenHandler := handler.NewTokenHandler(tokenUsecase, clientUsecase, clock)
			serverHandler := handler.NewServerHandler(accountInfoHandler, tokenHandler)
			presenter.RegisterHandlers(v1, serverHandler)
		}
	}

	return router, err
}
