package drop

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Read from the connection untill a new line is send
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the data read from the connection to the terminal
		fmt.Print("> ", string(data))
		os.Exit(0)
	}
}

// Broadcast your presence on the network
func broadcastPresence(tcpAddress string) {
	conn, err := net.Dial("udp", "255.255.255.255:6969")
	if err != nil {
		fmt.Println("Error setting up broadcast:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Broadcasting presence...")

	for {
		conn.Write([]byte(tcpAddress))
		time.Sleep(5 * time.Second)
	}
}

var receiveCommand = &cobra.Command{
	Use:   "receive",
	Short: "Receive file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Waiting for incoming files...")

		// Listen on any available port
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		fmt.Println("Listening on", listener.Addr().String())
		go broadcastPresence(listener.Addr().String())

		for {
			// Accept new connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
			}
			// Handle new connections in a Goroutine for concurrency
			go handleConnection(conn)
		}
	},
}
