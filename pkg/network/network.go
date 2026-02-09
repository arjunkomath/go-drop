package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// UDPPort is the port used for UDP communication
var UDPPort = 5050

// DevicePresenceMsg is a struct that represents the message sent by a device to
// broadcast its presence on the network
type DevicePresenceMsg struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// BroadcastPresence will broadcast device presence on the network
func BroadcastPresence(ctx context.Context, message DevicePresenceMsg) error {
	conn, err := net.Dial("udp", fmt.Sprintf("255.255.255.255:%d", UDPPort))
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}

	defer conn.Close()

	// Convert the struct to a JSON string
	jsonData, err := json.Marshal(message) // Marshaling converts struct to JSON
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if _, err := conn.Write(jsonData); err != nil {
				return fmt.Errorf("failed to write UDP: %w", err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// ParseDevicePresence parses the device presence message
func ParseDevicePresence(data []byte) (DevicePresenceMsg, error) {
	var message DevicePresenceMsg
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

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return nil, fmt.Errorf("unexpected address type: %T", conn.LocalAddr())
	}

	return localAddr.IP, nil
}
