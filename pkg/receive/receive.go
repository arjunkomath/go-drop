package receive

import (
	"bufio"
	"drop/pkg/device"
	"drop/pkg/network"
	"drop/pkg/utils"
	"drop/styles"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func waitForData() tea.Cmd {
	return func() tea.Msg {
		localIP, err := network.GetOutboundIP()
		if err != nil {
			return errorMsg(err)
		}

		name, err := device.GetName()
		if err != nil {
			return errorMsg(err)
		}

		listener, err := net.Listen("tcp", localIP.To4().String()+":0")
		if err != nil {
			return errorMsg(err)
		}
		defer listener.Close()

		presenseMsg := network.DevicePresenseMsg{
			Name:    name,
			Address: listener.Addr().String(),
		}

		go network.BroadcastPresence(presenseMsg)

		for {
			// Accept new connections
			conn, err := listener.Accept()
			if err != nil {
				return errorMsg(err)
			}

			defer conn.Close()

			for {
				// Read from the connection untill a new line is send
				data, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					return errorMsg(err)
				}

				return statusMsg(string(data))
			}
		}
	}
}

type errorMsg error
type statusMsg string

type model struct {
	spinner        spinner.Model
	secondsElapsed int
	waiting        bool
	errorMsg       string
	message        string
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner:        s,
		secondsElapsed: 0,
		waiting:        true,
		errorMsg:       "",
		message:        "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(utils.SecondTick(), m.spinner.Tick, tea.EnterAltScreen, waitForData())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case utils.TickMsg:
		m.secondsElapsed++
		return m, utils.SecondTick()

	case errorMsg:
		m.waiting = false
		m.errorMsg = msg.Error()
		return m, nil

	case statusMsg:
		m.waiting = false
		m.message = string(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := styles.HeaderStyle.Render("GoDrop")
	s += "\n\n"

	if m.waiting {
		s += fmt.Sprintf("%s Waiting for data... (%ds)\n\n", m.spinner.View(), m.secondsElapsed)
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

// Command for receiving data
var Command = &cobra.Command{
	Use:   "receive",
	Short: "Receive file",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}