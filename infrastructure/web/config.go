package web

import (
	"strings"

	"go-banking-api/pkg"
)

type Config struct {
	Host             string
	Port             string
	CorsAllowOrigins []string
}

func NewConfigWeb() *Config {
	return &Config{
		Host: pkg.GetEnvDefault("WEB_HOST", "0.0.0.0"),
		Port: pkg.GetEnvDefault("WEB_PORT", "8080"),
		CorsAllowOrigins: strings.Split(pkg.GetEnvDefault("CORS_ALLOW_ORIGINS",
			"http://localhost:8001"), ","),
	}
}
