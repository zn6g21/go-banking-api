package router

import (
	"github.com/gin-gonic/gin"

	"go-banking-api/adapter/controller/gin/handler"
	"go-banking-api/adapter/controller/gin/middleware"
)

/*
// Swaggerの設定をする
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
*/

func NewGinRouter(corsAllowOrigins []string) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	/*swagger, err := setupSwagger(router)
	if err != nil {
		logger.Warn(err.Error())
		return nil, err
	} */

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	// Healthチェック用のAPIです
	router.GET("/health", handler.Health)

	return router
}
