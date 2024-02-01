package main

import "github.com/alecthomas/kingpin"

var (
	app = kingpin.New("beautyf", "api with limited rate")

	runCmd    = app.Command("run", "run svc")
	rateLimit = runCmd.Flag("rateLimit", "rate limit for each endpoint").Default("10").Float64()
)
