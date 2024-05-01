package styles

import "github.com/charmbracelet/lipgloss"

// HeaderStyle is used as app header
var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(6).
	PaddingRight(6).
	Align(lipgloss.Center)

// DeviceNameStyle is used to style device names
var DeviceNameStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA"))

// SelectedDeviceStyle is used to highlight selected device cursor
var SelectedDeviceStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4"))
