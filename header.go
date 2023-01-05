package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var active = spinner.Spinner{
	Frames: []string{"∙", "●"},
	FPS:    time.Second / 2,
}

type header struct {
	spinner     spinner.Model
	title       string
	directory   string
	files       int
	records     int
	width       int
	showSpinner bool
}

func NewHeader(directory string, files int, records int) header {
	s := spinner.New()
	s.Spinner = active
	s.Style = lipgloss.NewStyle().Foreground(Theme().primary)
	return header{
		spinner:     s,
		title:       "DMARC Reports",
		directory:   directory,
		files:       files,
		records:     records,
		width:       80,
		showSpinner: false,
	}
}

func (h header) Init() tea.Cmd { return nil }

func (h header) Update(msg tea.Msg) (header, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.(type) {
	case ScanResultsMsg:
		return h, nil
	case ScanTriggerMsg, spinner.TickMsg:
		// start to tick (or keep on)
		h.spinner, cmd = h.spinner.Update(spinner.Tick())
		return h, cmd
	}
	return h, nil
}

func (h header) View() string {
	baseStyle := lipgloss.NewStyle().MarginLeft(1)
	s := baseStyle.
		Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "255"}).
		Bold(true).
		Render(h.title)

	if h.showSpinner {
		s += fmt.Sprintf(" %s\n", h.spinner.View())
	} else {
		s += "\n"
	}
	s += baseStyle.
		Bold(false).
		Foreground(lipgloss.Color("240")).
		Render(fmt.Sprintf("Directory: %s\nFiles: %d  Records: %d",
			h.directory,
			h.files,
			h.records)) + "\n"
	return s
}
