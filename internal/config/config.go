package config

import (
	logstd "log"
	"math/rand"
)

type Config struct {
	RateLimit float64
	Logger    *logstd.Logger
	Random    *rand.Rand
}
