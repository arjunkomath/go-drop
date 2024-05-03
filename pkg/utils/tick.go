package utils

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is a message that is sent every tick.
type TickMsg time.Time

// SecondTick returns a command that ticks every second.
func SecondTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
