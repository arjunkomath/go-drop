package drop

import (
	"bufio"
	"drop/pkg/device"
	"drop/pkg/network"
	"fmt"
	"log"
	"net"
	"os"

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

var receiveCommand = &cobra.Command{
	Use:   "receive",
	Short: "Receive file",
	Run: func(cmd *cobra.Command, args []string) {
		localIP, err := network.GetOutboundIP()
		if err != nil {
			fmt.Println("Error getting local IP:", err)
			return
		}

		name, err := device.GetName()
		if err != nil {
			fmt.Println("Error getting hostname:", err)
			return
		}

		fmt.Println("Waiting for incoming files...")

		// Listen on any available port
		listener, err := net.Listen("tcp", localIP.To4().String()+":0")
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		presenseMsg := network.DevicePresenseMsg{
			Name:    name,
			Address: listener.Addr().String(),
		}
		fmt.Println("Listening on", listener.Addr().String())

		go network.BroadcastPresence(presenseMsg)

		for {
			// Accept new connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
			}

			go handleConnection(conn)
		}
	},
}
