package tollboothfiber

import (
	tlimiter "github.com/didip/tollbooth/limiter"
	"github.com/eonianmonk/tollbooth_fasthttp"
	"github.com/gofiber/fiber/v2"
)

func LimitHandler(handler func(*fiber.Ctx) error, limiter *tlimiter.Limiter) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.Context() // fasthttp ctx
		httpError := tollbooth_fasthttp.LimitByRequest(limiter, ctx)

		if httpError != nil {
			c.Set("Content-Type", limiter.GetMessageContentType())
			return c.Status(httpError.StatusCode).JSON(httpError)
		}
		return handler(c)
	}
}

func LimitMiddleware(limiter *tlimiter.Limiter) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.Context() // fasthttp context
		httpError := tollbooth_fasthttp.LimitByRequest(limiter, ctx)
		if httpError != nil {
			c.Set("Content-Type", limiter.GetMessageContentType())
			return c.Status(httpError.StatusCode).JSON(httpError)
		}
		return c.Next()
	}
}
