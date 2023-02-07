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
	table    *table.Model
	viewer   xmlViewer
	help     help.Model
	results  FeedbackResults
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
	table := NewTable(results)
	table.Focus()
	return model{
		scanner:  scanner,
		header:   NewHeader(dir, results.Files(), results.Len()),
		table:    &table,
		viewer:   NewXMLViewer(),
		help:     h,
		results:  results,
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

func (m model) selected() *FeedbackResult {
	selected := m.table.Cursor()
	return m.results[selected]
}

func (m model) nextFocus() {
	if m.table.Focused() {
		m.table.Blur()
		m.viewer.Focus()
	} else {
		m.table.Focus()
		m.viewer.Blur()
	}
	// fmt.Println("table:", m.table.Focused())
	// fmt.Println("viewer:", m.viewer.Focused())
}

func (m model) resize(width, height int) {
	m.table.SetHeight(height - 8)
	m.header.width = width
	m.viewer.SetHeight(height - 7)
	m.viewer.SetWidth(width - m.table.Width())
}

func (m model) Show() tea.Msg {
	selected := m.selected()
	xmlMsg := ShowXMLRecordMsg(string(selected.XML))
	return tea.Msg(xmlMsg)
}

func (m model) Init() tea.Cmd {
	m.table.Focus()
	m.viewer.Blur()
	// display the xml of the selected line
	return m.Show
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var tbl table.Model
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.updating {
			// pass the tick to the header
			// it returns a new tick
			m.header, cmd = m.header.Update(msg)
			cmds = append(cmds, cmd)
		}
	case ShowXMLRecordMsg:
		m.viewer, cmd = m.viewer.Update(msg)
		cmds = append(cmds, cmd)
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
		// m.table.SetWidth(msg.Width - 2)
		m.table.SetHeight(msg.Height - 8)
		m.header.width = msg.Width
		m.viewer.SetHeight(msg.Height - 7)
		m.viewer.SetWidth(msg.Width - m.table.Width())
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.nextFocus()
		case "s":
			var msg ScanTriggerMsg
			return m, func() tea.Msg { return tea.Msg(msg) }
		default:
			if m.table.Focused() {
				tbl, cmd = m.table.Update(msg)
				m.table = &tbl
				selected := m.selected()
				xmlMsg := ShowXMLRecordMsg(string(selected.XML))
				cmd2 := func() tea.Msg { return tea.Msg(xmlMsg) }
				cmds = append(cmds, cmd, cmd2)
			} else {
				m.viewer, cmd = m.viewer.Update(msg)
				cmds = append(cmds, cmd)
			}

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
	if m.viewer.Width() < 25 {
		v += RenderTable(m.table) + "\n"
	} else {
		v += lipgloss.JoinHorizontal(lipgloss.Top, RenderTable(m.table), m.viewer.View()) + "\n"
	}
	v += " " + m.help.View(m.keys()) + "\n"
	return v
}
