package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// UDPPort is the port used for UDP communication
var UDPPort = 5050

// DevicePresenseMsg is a struct that represents the message sent by a device to
// broadcast its presence on the network
type DevicePresenseMsg struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// BroadcastPresence will broadcast device presence on the network
func BroadcastPresence(message DevicePresenseMsg) {
	conn, err := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", UDPPort))
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

	for {
		conn.Write([]byte(jsonData))
		time.Sleep(1 * time.Second)
	}
}

// ParseDevicePresence parses the device presence message
func ParseDevicePresence(data []byte) (DevicePresenseMsg, error) {
	var message DevicePresenseMsg
	err := json.Unmarshal(data, &message)
	if err != nil {
		return message, err
	}
	return message, nil
}

// GetOutboundIP gets the local IP address
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
