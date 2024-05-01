package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// DevicePresenseMsg is a struct that represents the message sent by a device to
// broadcast its presence on the network
type DevicePresenseMsg struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Broadcast device presence on the network
func BroadcastPresence(message DevicePresenseMsg) {
	conn, err := net.Dial("udp", "255.255.255.255:6969")
	if err != nil {
		fmt.Println("Error setting up broadcast:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Broadcasting presence...")

	// Convert the struct to a JSON string
	jsonData, err := json.Marshal(message) // Marshaling converts struct to JSON
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Convert the byte slice to a string and print
	fmt.Println("JSON string:", string(jsonData))

	for {
		conn.Write([]byte(jsonData))
		time.Sleep(5 * time.Second)
	}
}
