package main

import (
	"flag"
	"fmt"
	"os"
	"rabbinator/cmd"
)

var (
	channel = flag.String("channel", "", "RabbitMQ channel to consume, it could be arbitrary name.")
	config  = flag.String("config", "", "Config for the RabbitMQ queue")
)

func main() {

	flag.Parse()

	// Channel flag is required.
	if *channel == "" {
		flag.PrintDefaults()
		fmt.Println("Channel name is required")
		os.Exit(1)
	}

	// Initialize and set config.
	cmd.ConfigSetup(*channel, *config)

	// RabbitMQ connect.
	cmd.Connect()

}