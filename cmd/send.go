package drop

import (
	"drop/pkg/network"
	"drop/styles"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type errorMsg error
type statusMsg string

type deviceFound struct {
	tcpIP string
	name  string
}

var allDevices = make([]string, 0)

func sendTCPMessage(device deviceFound) tea.Cmd {
	return func() tea.Msg {
		tcpAddr, err := net.ResolveTCPAddr("tcp", string(device.tcpIP))
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

func searchForDevices() tea.Cmd {
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

			deviceMessage, err := network.ParseDevicePresence(buffer[:n])
			if err != nil {
				continue
			}

			tcpIP := deviceMessage.Address

			if slices.Contains(allDevices, tcpIP) {
				continue
			}

			allDevices = append(allDevices, tcpIP)

			return deviceFound{
				tcpIP: tcpIP,
				name:  deviceMessage.Name,
			}
		}
	}
}

type model struct {
	spinner        spinner.Model
	secondsElapsed int
	searching      bool

	devices []deviceFound
	cursor  int

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
		devices:   []deviceFound{},
		searching: true,

		cursor:         0,
		secondsElapsed: 0,

		errorMsg: "",
		message:  "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), m.spinner.Tick, searchForDevices(), tea.EnterAltScreen)
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

	case deviceFound:
		m.devices = append(m.devices, msg)
		return m, searchForDevices()

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.devices)-1 {
				m.cursor++
			}

		case "enter", " ":
			return m, sendTCPMessage(m.devices[m.cursor])
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

func (m model) View() string {
	s := styles.HeaderStyle.Render("GoDrop")
	s += "\n\n"

	if m.searching {
		s += fmt.Sprintf("%s Searching for devices... (%ds)\n\n", m.spinner.View(), m.secondsElapsed)

		for i, choice := range m.devices {
			cursor := " "
			if m.cursor == i {
				cursor = styles.SelectedDeviceStyle.Render(">")
			}

			s += fmt.Sprintf("%s %s (%s)\n", cursor, styles.DeviceNameStyle.Render(choice.tcpIP), choice.name)
		}
	}

	if m.errorMsg != "" {
		s += fmt.Sprintf("Error: %s\n", m.errorMsg)
	}

	if m.message != "" {
		s += fmt.Sprintf("%s\n", m.message)
	}

	s += "\n(press q to quit)\n"

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
