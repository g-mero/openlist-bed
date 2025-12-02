package middleware

import (
	"openlist-bed/pkg/R"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

func Auth() fiber.Handler {
	return func(c fiber.Ctx) error {
		if c.Get("API-KEY") != viper.GetString("api_key") {
			return R.ErrApiKeyInvalid
		}

		return c.Next()
	}
}
