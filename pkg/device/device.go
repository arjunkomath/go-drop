package device

import (
	"fmt"
	"os"
)

// GetName returns the name of the device
func GetName() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return "", err
	}

	return hostname, nil
}
