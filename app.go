package main

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

// var baseStyle = lipgloss.NewStyle()

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	table.KeyMap
	Scan key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Scan, k.Quit, k.LineUp, k.LineDown, k.PageDown, k.PageUp}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Scan, k.LineUp},
		{k.Quit, k.LineDown},
		{k.PageDown, k.PageUp},
		// {k.LineUp, k.LineUp, k.LineUp},
		// {k.LineUp, k.LineDown, k.PageDown, k.PageUp}, // second column
	}
}

// var keys = keyMap{
// 	Up: key.NewBinding(
// 		key.WithKeys("up"),
// 		key.WithHelp("↑", "move up"),
// 	),
// 	Down: key.NewBinding(
// 		key.WithKeys("down"),
// 		key.WithHelp("↓", "move down"),
// 	),
// 	Quit: key.NewBinding(
// 		key.WithKeys("q", "esc", "ctrl+c"),
// 		key.WithHelp("q", "quit"),
// 	),
// }

type model struct {
	scanner  scanner
	header   header
	table    table.Model
	help     help.Model
	updating bool
}

func NewModel(directory string) model {
	dir, err := filepath.Abs(directory)
	if err != nil {
		dir = directory
	}
	scanner := NewScanner(dir)
	results := scanner.rawScan()

	h := help.New()
	// h.ShowAll = false
	return model{
		scanner:  scanner,
		header:   NewHeader(dir, results.Files(), results.Len()),
		table:    NewTable(results),
		help:     h,
		updating: false,
	}
}

func (m model) keys() keyMap {
	return keyMap{
		KeyMap: m.table.KeyMap,
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Scan: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "scan directory"),
		)}
}

func (m model) Init() tea.Cmd {
	m.table.Focus()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		// fmt.Println(m.scanner.busy)
		if m.updating {
			// pass the tick to the header
			// it returns a new tick
			m.header, cmd = m.header.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ScanResultsMsg:
		// receive results from scanner
		_, rows := toTable(FeedbackResults(msg))
		m.table.SetRows(rows)
		m.updating = false
		m.header.showSpinner = false
	case ScanTriggerMsg:
		m.updating = true
		m.header.showSpinner = true
		// update the header
		m.header, cmd = m.header.Update(msg)
		cmds = append(cmds, cmd)
		// update the scanner
		m.scanner, cmd = m.scanner.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width - 2)
		m.table.SetHeight(msg.Height - 8)
		m.header.width = msg.Width
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.table.MoveUp(1)
		case "down", "j":
			m.table.MoveDown(1)
		case "shift+up":
			m.table.MoveUp(3)
		case "shift+down":
			m.table.MoveDown(3)
		case "f", "pgdown":
			m.table.GotoBottom()
		case "b", "pgup":
			m.table.GotoTop()
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "s":
			var msg ScanTriggerMsg
			return m, func() tea.Msg { return tea.Msg(msg) }
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)

		}
	case tea.MouseEvent:
		fmt.Println(msg)
		switch msg.Type {
		case tea.MouseWheelDown:
			m.table.MoveDown(msg.Y)
		case tea.MouseWheelUp:
			m.table.MoveUp(msg.Y)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	v := m.header.View()
	v += baseStyle.Render(m.table.View()) + "\n"
	v += " " + m.help.View(m.keys()) + "\n"
	return v
}
