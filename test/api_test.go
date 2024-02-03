package test

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/eonianmonk/dummy-beutifier/internal/config"
	backend "github.com/eonianmonk/dummy-beutifier/internal/http"
	"github.com/gofiber/fiber/v2"
)

func client() http.Client {
	return *http.DefaultClient
}

var (
	address = "localhost:8080"

	helloEndpoint   = fmt.Sprintf("http://%s/hello", address)
	jsonEndpoint    = fmt.Sprintf("http://%s/beautify/json", address)
	jsonAPIEndpoint = fmt.Sprintf("http://%s/beautify/jsonapi", address)
)

var serverActive = false

func startServer(t *testing.T) {
	if serverActive {
		return
	}
	serverActive = true
	logger := log.Default()
	logger.SetOutput(os.Stderr)
	cfg := config.Config{
		RateLimit: 10,
		Logger:    logger,
		Random:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Endpoint:  address,
	}
	t.Log("starting backend...")
	err := backend.StartFiber(cfg)
	panic(err.Error())
}

func TestAPI(t *testing.T) {
	go startServer(t)
	time.Sleep(time.Second * 1)

	t.Run("basic-valid-requests", func(t *testing.T) {
		go startServer(t)
		time.Sleep(time.Second * 1)

		validJSON := `{"username":"konishe","id":23152,"links":["2.com","1.com","bb.com"]}`
		validJSONAPI := `{"data":{"type":"post","id":"123","attributes":{" created_at":-61758633600,"author":"Alan Donovan","body":"golang.org/x/tools/cmd/deadcode@latest","title":"Finding unreachable functions with deadcode"},"relationships":{"comments":{"data":[{"type":"comment","id":"0"},{"type":"comment","id":"1"}]}}},"included":[{"type":"comment","id":"0","attributes":{" created_at":-61758633600,"author":"a0@mail.com","body":"Wow! useful"}},{"type":"comment","id":"1","attributes":{" created_at":-61758630000,"author":"a1@mail.com","body":"Wow! useful x2"}}]}`

		// default indent is 2 spaces - currently hardcoded
		//
		validJSONResult := `{
  "username": "konishe",
  "id": 23152,
  "links": [
    "2.com",
    "1.com",
    "bb.com"
  ]
}`
		validJSONAPIResult := `{
  "data": {
    "type": "post",
    "id": "123",
    "attributes": {
      " created_at": -61758633600,
      "author": "Alan Donovan",
      "body": "golang.org/x/tools/cmd/deadcode@latest",
      "title": "Finding unreachable functions with deadcode"
    },
    "relationships": {
      "comments": {
        "data": [
          {
            "type": "comment",
            "id": "0"
          },
          {
            "type": "comment",
            "id": "1"
          }
        ]
      }
    }
  },
  "included": [
    {
      "type": "comment",
      "id": "0",
      "attributes": {
        " created_at": -61758633600,
        "author": "a0@mail.com",
        "body": "Wow! useful"
      }
    },
    {
      "type": "comment",
      "id": "1",
      "attributes": {
        " created_at": -61758630000,
        "author": "a1@mail.com",
        "body": "Wow! useful x2"
      }
    }
  ]
}`

		cli := client()
		helloResultPattern := regexp.MustCompile(`Hello after sleeping for \d+\.\d{3}`)

		resp, err := cli.Get(helloEndpoint)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("failed to fetch %s", helloEndpoint)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("failed to read body: ", err.Error())
		}
		if !helloResultPattern.MatchString(string(body)) {
			t.Fatal("unexpected hello endpoint response pattern: ", string(body))
		}

		req1Body := strings.NewReader(validJSON)
		resp1, err := cli.Post(jsonEndpoint, "application/json", req1Body)
		if err != nil {
			t.Fatalf("failed to post to %s endpoint: %s", jsonEndpoint, err.Error())
		}
		defer resp1.Body.Close()
		body, err = io.ReadAll(resp1.Body)
		if err != nil {
			t.Fatal("failed to read json resp body: ", err.Error())
		}
		if string(body) != validJSONResult {
			t.Fatal("unexpected result on beautified json")
		}

		req2Body := strings.NewReader(validJSONAPI)
		resp2, err := cli.Post(jsonAPIEndpoint, "application/json", req2Body)
		if err != nil {
			t.Fatalf("failed to post to %s: %s", jsonAPIEndpoint, err.Error())
		}
		defer resp2.Body.Close()
		body, err = io.ReadAll(resp2.Body)
		if err != nil {
			t.Fatal("failed to read json resp body: ", err.Error())
		}
		if string(body) != validJSONAPIResult {
			fmt.Println(string(body))
			fmt.Println(validJSONAPIResult)
			t.Fatal("unexpected result on beautified json")
		}

	})

	t.Run("invalid-json-requests", func(t *testing.T) {
		go startServer(t)
		time.Sleep(time.Second * 1)
		cli := client()

		reqFn := func(endp string, expectedCode int, req string) {
			reqBody := strings.NewReader(req)
			resp, err := cli.Post(endp, "application/json", reqBody)
			if err != nil {
				t.Fatalf("failed to post to %s endpoint: %s", jsonEndpoint, err.Error())
			}
			defer resp.Body.Close()
			if resp.StatusCode != expectedCode {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("unexpected endpoint response code on invalid json at %s\nReceived code: %d, wanted: %d\nResponse body: %s",
					endp, resp.StatusCode, expectedCode, string(body))
			}
		}

		invalidJSON := "{"
		invalidJSONAPI := "{"
		invalidJSONAPISchema := `{"weird_key":"weird_value"}`

		reqFn(jsonEndpoint, fiber.StatusBadRequest, invalidJSON)

		// invalid json
		reqFn(jsonAPIEndpoint, fiber.StatusBadRequest, invalidJSONAPI)

		// invalid jsonapi scheme
		reqFn(jsonAPIEndpoint, fiber.StatusBadRequest, invalidJSONAPISchema)
	})
}
