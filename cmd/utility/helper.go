package utility

import "log"

// Handling errors.
func ErrorLog(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}