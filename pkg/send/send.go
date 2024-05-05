package send

import (
	"drop/pkg/network"
	"drop/styles"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/textarea"
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

func sendTCPMessage(device deviceFound, message string) tea.Cmd {
	return func() tea.Msg {
		tcpAddr, err := net.ResolveTCPAddr("tcp", string(device.tcpIP))
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			return errorMsg(err)
		}

		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			return errorMsg(err)
		}

		defer conn.Close()
		return statusMsg("Message sent")
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

	devices []deviceFound
	cursor  int
	input   string

	sending  bool
	textarea textarea.Model

	errorMsg string
	message  string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textarea.New()
	ti.Placeholder = "Hello world..."
	ti.Focus()

	sw := stopwatch.NewWithInterval(time.Second)

	return model{
		spinner:   s,
		stopwatch: sw,

		searching: true,
		devices:   []deviceFound{},
		cursor:    0,

		sending:  false,
		textarea: ti,

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
			m.sending = true
			return m, textarea.Blink

		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}

		case tea.KeyCtrlS:
			m.sending = false
			return m, sendTCPMessage(m.devices[m.cursor], m.textarea.Value())

		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := styles.HeaderStyle.Render("GoDrop")
	s += "\n\n"

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

	if m.sending {
		s += fmt.Sprintf(
			"Enter message.\n\n%s\n",
			m.textarea.View(),
		)
	}

	if m.errorMsg != "" {
		s += fmt.Sprintf("Error: %s\n", m.errorMsg)
	}

	if m.message != "" {
		s += fmt.Sprintf("%s\n", m.message)
	}

	if m.sending {
		s += "\n(press ctrl+s to send, ctrl+c to quit)\n"
	} else {
		s += "\n(press ctrl+c to quit)\n"
	}

	return s
}

// Command used for sending data
var Command = &cobra.Command{
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
