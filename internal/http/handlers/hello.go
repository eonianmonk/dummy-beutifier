package handlers

import (
	"fmt"
	"time"

	"github.com/eonianmonk/dummy-beutifier/internal/http/middleware"
	"github.com/gofiber/fiber/v2"
)

// sleeps for 500ms-1.5s and sends hello world
func Hello(c *fiber.Ctx) error {
	rand := middleware.GetRand(c)
	sleepTime := (rand.Int() % 1000) + 500
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	resp := fmt.Sprintf("Hello after sleeping for %.3f", float64(sleepTime)/1000.0)
	return c.Status(200).SendString(resp)
}
