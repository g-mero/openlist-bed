package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/utils/v2"
)

func Cache() fiber.Handler {
	return cache.New(cache.Config{
		KeyGenerator: func(c fiber.Ctx) string {
			return utils.CopyString(c.Path() +
				c.Get("Accept") +
				c.Query("webp"))
		},
		MaxBytes:   1024 * 1024 * 10, // 10 MB
		Expiration: time.Hour,
	})
}
