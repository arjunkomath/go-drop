package drop

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type errorMsg error
type statusMsg string

func startSearching() tea.Cmd {
	return func() tea.Msg {
		addr, err := net.ResolveUDPAddr("udp", ":6969")
		if err != nil {
			return errorMsg(err)
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return errorMsg(err)
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue
			}

			tcpAddr, err := net.ResolveTCPAddr("tcp", string(buffer[:n]))
			// Connect to the address with tcp
			conn, err := net.DialTCP("tcp", nil, tcpAddr)

			if err != nil {
				return errorMsg(err)
			}

			// Send a message to the server
			_, err = conn.Write([]byte("Hello TCP Server\n"))
			if err != nil {
				return errorMsg(err)
			}

			defer conn.Close()
			return statusMsg("Message sent")
		}
	}
}

type model struct {
	spinner        spinner.Model
	secondsElapsed int
	searching      bool

	devices  []string
	cursor   int
	selected map[int]struct{}

	errorMsg string
	message  string
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner:   s,
		devices:   []string{},
		searching: true,

		cursor:         0,
		selected:       make(map[int]struct{}),
		secondsElapsed: 0,

		errorMsg: "",
		message:  "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), m.spinner.Tick, startSearching(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.secondsElapsed++
		return m, tick()

	case errorMsg:
		m.searching = false
		m.errorMsg = msg.Error()
		return m, tea.Quit

	case statusMsg:
		m.searching = false
		m.message = string(msg)
		return m, nil

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

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(4).
	Width(32)

func (m model) View() string {
	s := headerStyle.Render("GoDrop")
	s += "\n\n"

	if m.searching {
		s += fmt.Sprintf("%s Searching for devices... (%d)\n\n", m.spinner.View(), m.secondsElapsed)
	}

	if m.errorMsg != "" {
		s += fmt.Sprintf("Error: %s\n\n", m.errorMsg)
	}

	if m.message != "" {
		s += fmt.Sprintf("Status: %s\n\n", m.message)
	}

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

	s += "\nPress q to quit.\n"

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
