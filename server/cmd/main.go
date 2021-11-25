package main

import (
	"fmt"
	"os"

	"gitlab.com/mirukakoro/sisikyo/server"
)

func main() {
	err := server.Main()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
