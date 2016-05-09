package logger

import "log"

// Process check error and log it if requires
func Process(err error, text string) {
	if err != nil {
		log.Fatal(text+" ", err)
	}
}
