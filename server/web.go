package server

import (
	"flag"
	"fmt"
)

var webPort int
var webHost string

func init() {
	flag.IntVar(&webPort, "port", 8080, "port to bind to")
	flag.StringVar(&webHost, "host", "", "host to bind to")
}

func setupWeb() (string, error) {
	return fmt.Sprintf("%s:%d", webHost, webPort), nil
}
