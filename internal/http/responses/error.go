package responses

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(c *fiber.Ctx, message string, status int) error {
	errr := ErrorResponse{Error: message}
	c.Status(status).JSON(errr)
	return nil
}
