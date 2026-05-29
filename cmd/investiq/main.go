package main

import (
	"fmt"
	"os"

	"github.com/dkstm95/investiq/internal/app"
)

var version = "0.1.0"

func main() {
	application := app.New(version)
	if err := application.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
