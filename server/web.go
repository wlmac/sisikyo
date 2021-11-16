package server

import (
	"flag"
	"fmt"
	"time"
)

var webPort int
var webHost string
var cacheTimeout time.Duration

func init() {
	flag.IntVar(&webPort, "port", 8080, "port to bind to")
	flag.StringVar(&webHost, "host", "", "host to bind to")
	flag.DurationVar(&cacheTimeout, "cache", time.Hour, "cache timeout")
}

func setupWeb() (string, error) {
	return fmt.Sprintf("%s:%d", webHost, webPort), nil
}
