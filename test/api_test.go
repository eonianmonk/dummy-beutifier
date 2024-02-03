package test

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/eonianmonk/dummy-beutifier/internal/config"
	backend "github.com/eonianmonk/dummy-beutifier/internal/http"
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

func startServer(t *testing.T) {
	logger := log.Default()
	logger.SetOutput(nil)
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
		validJSON := `{"username":"konishe","id":23152,"links":["2.com","1.com","bb.com"]}`
		validJSONAPI := `{"links":{"self":"http://example.com/v1/api/resource/123"},"data":{"id":"123","type":"your_resources"},"included":[{"id":"456","type":"related_resources"}]}`

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
  "links": {
    "self": "http://example.com/v1/api/resource/123"
  },
  "data": {
    "id": "123",
    "type": "your_resources"
  },
  "included": [
    {
      "id": "456",
      "type": "related_resources"
    }
  ]
}`
		helloResultPattern := regexp.MustCompile(`Hello after sleeping for \d+\.\d{3}`)

		cli := client()

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
}
