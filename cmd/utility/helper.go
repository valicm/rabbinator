package utility

import (
	"log"
)

// Definition for queue item status.
type QueueStatus struct {
	Success string `yaml:"default: success"`
	Reject string `yaml:"default: reject"`
	Retry string `yaml:"default: retry"`
	Unknown string `yaml:"default: unknown"`
}

// Handling errors.
func ErrorLog(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
