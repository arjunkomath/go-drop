package drop

import (
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	_, err = conn.Write([]byte("Hello from TCP Server\n"))
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

type model struct {
	devices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		devices:  []string{},
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	go listenForBroadcast()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.devices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := fmt.Sprintf("Searching for devices... (%d)\n\n", len(m.devices))

	// Iterate over our choices
	for i, choice := range m.devices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

var sendCommand = &cobra.Command{
	Use:   "send",
	Short: "Send file",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}
