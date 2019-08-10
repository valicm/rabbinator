package main

import (
	"flag"
	"fmt"
	"os"
	"rabbinator/cmd"
)

var (
	consumer = flag.String("consumer", "", "Consumer tag, should be unique. Used for distinction between multiple consumers.")
	config  = flag.String("config", "", "Optional. Declare specific directory where config files are located. Etc. /var/www/my_directory")
)

func main() {

	flag.Parse()

	// Consumer flag is required.
	if *consumer == "" {
		flag.PrintDefaults()
		fmt.Println("Consumer flag is required. It is used to distinct multiple consumers for same queue, and utilizes yaml configuration with same naming")
		os.Exit(1)
	}

	// Initialize configuration setup and RabbitMQ connection.
	cmd.Initialize(*consumer, *config)

}