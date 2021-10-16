package main

import (
	"fmt"
	"gitlab.com/mirukakoro/sisikyo/server"
	"os"
)

func main() {
	err := server.Main()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
}
