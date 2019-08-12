package utility

import (
	"fmt"
	"os"
)

// Helper for printing errors and exiting the program.
// Use for initialization task before we enter forever loop
// for processing queues. Setup config, connecting to RabbitMQ.
func InitErrorHandler(message string, err error) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}

// User related errors during configuration, etc..
// Not really any dev errors probably. Just unsupported choices.
func InputErrorHandler(message string) {
	fmt.Print(message + "\n")
	os.Exit(1)
}
