package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/spf13/viper"
)

func Cors() fiber.Handler {
	var allowOrigins = []string{"*"}
	allowHost := viper.GetString("allow_origins")
	if allowHost != "" {
		allowOrigins = strings.Split(allowHost, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins: allowOrigins,
	})
}
