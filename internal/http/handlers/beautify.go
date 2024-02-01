package handlers

import (
	"encoding/json"
	"time"

	"github.com/eonianmonk/dummy-beutifier/internal/http/middleware"
	"github.com/eonianmonk/dummy-beutifier/internal/http/responses"
	"github.com/gofiber/fiber/v2"
)

// function beautifies JSON
func BeautifyJSON(c *fiber.Ctx) error {
	log := middleware.GetLogger(c)

	var incomingData interface{}
	err := json.Unmarshal(c.Body(), incomingData)
	if err != nil { // invalid json
		msg := "invalid json schema request"
		log.Printf("%s: %s\n", msg, string(c.Body()))
		return responses.SendErrorResponse(c, msg, fiber.StatusBadRequest)
	}
	beautifiedJSON, err := json.MarshalIndent(incomingData, "", "  ")
	if err != nil {
		msg := "failed to marshalIndent"
		log.Printf("%s: %s", msg, err.Error())
		return responses.SendErrorResponse(c, msg, fiber.StatusInternalServerError)
	}

	rand := middleware.GetRand(c)
	sleepTime := (rand.Int() % 1000) + 500
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	return c.Status(fiber.StatusOK).SendString(string(beautifiedJSON))
}
