package device

import (
	"log"
	"os"
)

// GetName returns the name of the device
func GetName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Error getting hostname:", err)
	}

	return hostname, nil
}
