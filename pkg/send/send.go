package send

import (
	"drop/pkg/network"
	"drop/styles"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
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

func sendFile(device deviceFound, file string) tea.Cmd {
	return func() tea.Msg {
		tcpAddr, err := net.ResolveTCPAddr("tcp", string(device.tcpIP))
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			return errorMsg(err)
		}

		// Open the file to send
		file, err := os.Open(file) // Make sure the file exists
		if err != nil {
			return errorMsg(err)
		}
		defer file.Close()

		_, err = io.Copy(conn, file) // Send file content to the server
		if err != nil {
			return errorMsg(err)
		}

		defer conn.Close()
		return statusMsg("File sent")
	}
}

func searchForDevices() tea.Cmd {
	return func() tea.Msg {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", network.UDPPort))
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
	stopwatch stopwatch.Model
	spinner   spinner.Model
	searching bool
	file      string

	devices []deviceFound
	cursor  int
	input   string

	errorMsg string
	message  string
}

func initialModel(file string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	sw := stopwatch.NewWithInterval(time.Second)

	return model{
		spinner:   s,
		stopwatch: sw,
		file:      file,
		searching: true,

		devices: []deviceFound{},
		cursor:  0,

		errorMsg: "",
		message:  "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.stopwatch.Init(), searchForDevices(), tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case errorMsg:
		m.searching = false
		m.errorMsg = msg.Error()
		return m, nil

	case statusMsg:
		m.searching = false
		m.message = string(msg)
		return m, nil

	case deviceFound:
		m.devices = append(m.devices, msg)
		return m, searchForDevices()

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown:
			if m.cursor < len(m.devices)-1 {
				m.cursor++
			}

		case tea.KeyEnter:
			m.searching = false
			return m, sendFile(m.devices[m.cursor], m.file)

		default:

		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := styles.HeaderStyle.Render("GoDrop")
	s += "\n\n"

	s += fmt.Sprintf("Sending file: %s\n\n", m.file)

	if m.searching {
		s += fmt.Sprintf("%s Searching for devices... (%s)\n\n", m.spinner.View(), m.stopwatch.View())

		for i, choice := range m.devices {
			cursor := " "
			if m.cursor == i {
				cursor = styles.SelectedDeviceStyle.Render(">")
			}

			s += fmt.Sprintf("%s %s\t%s\n", cursor, styles.DeviceNameStyle.Render(choice.name), choice.tcpIP)
		}
	}

	if m.errorMsg != "" {
		s += fmt.Sprintf("Error: %s\n", m.errorMsg)
	}

	if m.message != "" {
		s += fmt.Sprintf("%s\n", m.message)
	}

	s += "\n(press ctrl+c to quit)\n"

	return s
}

// Command used for sending data
var Command = &cobra.Command{
	Use:   "send",
	Short: "Send file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		p := tea.NewProgram(initialModel(filePath))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}
