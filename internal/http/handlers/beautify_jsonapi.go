package handlers

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/eonianmonk/dummy-beutifier/internal/http/middleware"
	"github.com/eonianmonk/dummy-beutifier/internal/http/requests"
	"github.com/eonianmonk/dummy-beutifier/internal/http/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/jsonapi"
	"gitlab.com/c0b/go-ordered-json"
)

// function beautifies JSON
func BeautifyJSONAPI(c *fiber.Ctx) error {
	log := middleware.GetLogger(c)

	reqReader := bytes.NewReader(c.Body())
	// we don't know the incoming request schema
	// and UnmarshalPayload does not read to interface{}
	empty := requests.Article{}
	// verifying jsonapi schema
	err := jsonapi.UnmarshalPayload(reqReader, &empty)
	if err != nil {
		msg := "invalid schema"
		return responses.SendErrorResponse(c, msg, fiber.StatusBadRequest)
	}
	store := ordered.NewOrderedMap()
	_ = json.Unmarshal(c.Body(), &store)
	beautifiedJSON, err := json.MarshalIndent(store, "", "  ")
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
