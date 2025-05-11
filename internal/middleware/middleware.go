// internal/middleware/middleware.go
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/swiftcal/pkg/logger"
	"go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request
		duration := time.Since(start)

		logger.GetLogger().Info("HTTP request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}
