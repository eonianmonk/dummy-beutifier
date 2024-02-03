package main

import (
	stdlog "log"
	"math/rand"
	"os"
	"time"

	"github.com/eonianmonk/dummy-beutifier/internal/config"
	"github.com/eonianmonk/dummy-beutifier/internal/http"
	"github.com/pkg/errors"
)

func main() {
	Run(os.Args)
}

func Run(args []string) {
	randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))
	log := stdlog.Default()
	defer func() {
		if rvr := recover(); rvr != nil {
			log.Fatal(rvr, "-> app panicked")
		}
	}()

	cmd, err := app.Parse(args)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Failed to parse cli command"))
	}
	switch cmd {
	case runCmd.FullCommand():
		cfg := config.Config{
			RateLimit: *rateLimit,
			Logger:    log,
			Random:    randSrc,
			Endpoint:  "localhost:8080",
		}
		log.Println("Starting server")
		http.StartFiber(cfg)
	default:
		log.Panicf("unknown command")
	}
}
