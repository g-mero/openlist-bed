package middleware

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v3"
)

func ZapLogger(zapLogger *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Set variables
		var (
			once       sync.Once
			errHandler fiber.ErrorHandler
			errPadding = 15
		)

		// Set error handler once
		once.Do(func() {
			// get the longest possible path
			stack := c.App().Stack()
			for m := range stack {
				for r := range stack[m] {
					if len(stack[m][r].Path) > errPadding {
						errPadding = len(stack[m][r].Path)
					}
				}
			}
			// override error handler
			errHandler = c.App().Config().ErrorHandler
		})

		// Handle request, store err for logging
		chainErr := c.Next()

		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		latency := time.Since(start).Milliseconds()

		zapProps := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Int64("latency_ms", latency),
			zap.String("ip", c.IP()),
		}

		if chainErr != nil {
			zapProps = append(zapProps, zap.String("error", chainErr.Error()))
		}

		if v, ok := c.Locals("user_id").(int); ok {
			zapProps = append(zapProps, zap.Int("user_id", v))
		}

		zapLogger.Info("request", zapProps...)

		return nil
	}
}
