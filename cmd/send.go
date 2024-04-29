package drop

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

func sendTCPMessage(addr *net.TCPAddr) {
	// Connect to the address with tcp
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Send a message to the server
	_, err = conn.Write([]byte("Hello TCP Server\n"))
	fmt.Println("send...")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Message sent")
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
	for {
		n, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		fmt.Printf("Received message '%s' from %s\n", src.String(), string(buffer[:n]))

		tcpAddr, err := net.ResolveTCPAddr("tcp", string(buffer[:n]))
		// Connect to the address with tcp
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Send a message to the server
		_, err = conn.Write([]byte("Hello TCP Server\n"))
		fmt.Println("send...")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Message sent")
		defer conn.Close()
		break
	}
}

var sendCommand = &cobra.Command{
	Use:   "send",
	Short: "Send file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Searching for devices...")
		listenForBroadcast()
	},
}
