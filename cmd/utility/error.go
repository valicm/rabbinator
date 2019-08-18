package utility

import (
	"fmt"
	"os"
)

// InitErrorHandler print errors and exit the program.
// Use for initialization task before we enter forever loop
// for processing queues. Setup config, connecting to RabbitMQ.
func InitErrorHandler(message string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}

// inputErrorHandler handle user related errors during configuration, etc..
// Not really any dev errors probably. Just unsupported choices.
func inputErrorHandler(message string) {
	fmt.Print(message + "\n")
	os.Exit(1)
}
