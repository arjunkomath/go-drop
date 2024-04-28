package main

import (
	"fmt"
	"net"
	"slices"
	"time"
)

// Broadcast your presence on the network
func broadcastPresence() {
	conn, err := net.Dial("udp", "255.255.255.255:6969")
	if err != nil {
		fmt.Println("Error setting up broadcast:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Broadcasting presence...")

	for {
		conn.Write([]byte("Hello, I'm here!"))
		time.Sleep(5 * time.Second)
	}
}

// Listen for other devices on the network
func listenForBroadcast() {
	addr, err := net.ResolveUDPAddr("udp", ":6969")
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error setting up UDP listener:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	foundDevices := make([]string, 0)

	for {
		_, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		if !slices.Contains(foundDevices, src.String()) {
			foundDevices = append(foundDevices, src.String())
			fmt.Println("Found new device:", src.String())
		}

		fmt.Println("Found devices:", foundDevices)
	}
}

func main() {
	// broadcastPresence()

	go broadcastPresence()
	listenForBroadcast()
}
