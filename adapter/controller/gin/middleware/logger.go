package middleware

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"go-banking-api/pkg/logger"
)

func GinZap() gin.HandlerFunc {
	return ginzap.Ginzap(logger.ZapLogger, time.RFC3339, true)
}

func RecoveryWithZap() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(logger.ZapLogger, true)
}
