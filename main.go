package main

import (
	"./config"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	var opts config.Opts

	p := flags.NewParser(&opts, 0)

	_, err := p.Parse()

	if err != nil {
		fmt.Printf("Failed to parse args: %v\n", err)
		os.Exit(1)
	}
}
