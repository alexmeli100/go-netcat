package config

import "fmt"

type Opts struct {
	Listen            bool   `short:"l" long:"listen" description:"Listen for incoming connections" default:"false"`
	Command           bool   `short:"c" long:"command" description:"Initialize a command shell" default:"false"`
	Execute           string `short:"e" long:"execute" description:"Execute the given file" default:""`
	UploadDestination string `short:"u" long:"upload" description:"Upload destination for file" default:""`
	Addr              string `short:"a" long:"address" description:"The address to listen to." default:""`
	Port              int    `short:"p" long:"port" description:"The port to listen on" default:"8080"`
}

func usage() {
	fmt.Println("Go Netcat tool")
}
