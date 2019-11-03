package main

import (
	"fmt"
	"github.com/alexmeli100/go-netcat/config"
	"github.com/alexmeli100/go-netcat/driver"
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

	d := driver.Driver{Params: &opts}

	d.Run()
}
