package middleware

import (
	logstd "log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

var (
	logKey  string = "logger"
	randKey        = "rand"
)

func SetLogger(log *logstd.Logger) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(logKey, log)
		return c.Next()
	}
}

func GetLogger(c *fiber.Ctx) *logstd.Logger {
	return c.Locals(logKey).(*logstd.Logger)
}

func SetRand(rand *rand.Rand) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Locals(randKey, rand)
		return c.Next()
	}
}

func GetRand(c *fiber.Ctx) *rand.Rand {
	return c.Locals(randKey).(*rand.Rand)
}
