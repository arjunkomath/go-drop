package main

import drop "drop/cmd"

func main() {
	drop.Execute()
}

// package main

// import (
// 	"fmt"
// 	"net"
// 	"os"
// 	"slices"
// 	"time"
// )

// // Broadcast your presence on the network
// func broadcastPresence() {
// 	conn, err := net.Dial("udp", "255.255.255.255:6969")
// 	if err != nil {
// 		fmt.Println("Error setting up broadcast:", err)
// 		return
// 	}

// 	defer conn.Close()

// 	fmt.Println("Broadcasting presence...")

// 	for {
// 		conn.Write([]byte("Hello, I'm here!"))
// 		time.Sleep(5 * time.Second)
// 	}
// }

// func sendTCPMessage(addr *net.TCPAddr) {
// 	// Connect to the address with tcp
// 	conn, err := net.DialTCP("tcp", nil, addr)

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	// Send a message to the server
// 	_, err = conn.Write([]byte("Hello TCP Server\n"))
// 	fmt.Println("send...")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("Message sent")
// }

// // Listen for other devices on the network
// func listenForBroadcast() {
// 	addr, err := net.ResolveUDPAddr("udp", ":6969")
// 	if err != nil {
// 		fmt.Println("Error resolving UDP address:", err)
// 		return
// 	}

// 	conn, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		fmt.Println("Error setting up UDP listener:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	buffer := make([]byte, 1024)
// 	foundDevices := make([]string, 0)

// 	for {
// 		_, src, err := conn.ReadFromUDP(buffer)
// 		if err != nil {
// 			fmt.Println("Error reading from UDP:", err)
// 			continue
// 		}

// 		if !slices.Contains(foundDevices, src.String()) {
// 			foundDevices = append(foundDevices, src.String())
// 			fmt.Println("Found new device:", src.String())
// 		}

// 		if len(foundDevices) > 0 {
// 			fmt.Println("Found all devices, sending file...")

// 			tcpAddress := &net.TCPAddr{
// 				IP:   src.IP,   // Use the same IP address
// 				Port: src.Port, // Use the same port number
// 			}
// 			sendTCPMessage(tcpAddress)

// 			break
// 		}

// 		fmt.Println("Found devices:", foundDevices)
// 	}
// }

// func main() {
// 	broadcastPresence()

// 	// go broadcastPresence()
// 	// listenForBroadcast()
// }
