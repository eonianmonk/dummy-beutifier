package http

import (
	"github.com/didip/tollbooth"
	"github.com/eonianmonk/dummy-beutifier/internal/config"
	"github.com/eonianmonk/dummy-beutifier/internal/http/handlers"
	"github.com/eonianmonk/dummy-beutifier/internal/http/middleware"
	tollboothfiber "github.com/eonianmonk/dummy-beutifier/pkg/tollbooth_fiber"
	"github.com/gofiber/fiber/v2"
)

func StartFiber(cfg config.Config) error {
	app := fiber.New()

	limiter := tollbooth.NewLimiter(float64(cfg.RateLimit), nil)
	limiter.SetStatusCode(fiber.StatusTooManyRequests)
	app.Use(
		middleware.SetLogger(cfg.Logger),
		middleware.SetRand(cfg.Random),
		tollboothfiber.LimitMiddleware(limiter),
	)
	setupRoutes(app)

	return app.Listen(cfg.Endpoint)
}

func setupRoutes(app *fiber.App) {
	app.Get("/hello", handlers.Hello)

	beutifyApi := app.Group("/beautify")
	beutifyApi.Post("/json", handlers.BeautifyJSON)
	beutifyApi.Post("/jsonapi", handlers.BeautifyJSONAPI)
}
