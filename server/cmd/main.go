package main

import (
	"fmt"
	"os"

	"github.com/foolin/goview"
	"gitlab.com/mirukakoro/sisikyo/server"
)

func main() {
	fmt.Println(goview.DefaultConfig)
	err := server.Main()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
