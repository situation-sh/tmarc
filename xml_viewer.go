package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ShowXMLRecordMsg string

// component to read xml record
type xmlViewer struct {
	textarea  *textarea.Model
	highlight bool
}

func NewXMLViewer() xmlViewer {
	focused, blurred := textarea.DefaultStyles()
	focused.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(Theme().primary)
	blurred.LineNumber = lipgloss.NewStyle().Faint(true)
	blurred.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
	focused.LineNumber = lipgloss.NewStyle().Faint(true)

	t := textarea.New()
	t.BlurredStyle = blurred
	t.FocusedStyle = focused
	return xmlViewer{textarea: &t, highlight: highlightXML}
}

func (x xmlViewer) Init() tea.Cmd {
	return nil
}

func (x xmlViewer) Update(msg tea.Msg) (xmlViewer, tea.Cmd) {
	var cmd tea.Cmd
	var t textarea.Model
	switch m := msg.(type) {
	case ShowXMLRecordMsg:
		x.textarea.SetValue(string(m))
		if x.highlight {
			x.Highlight()
		}
	}

	t, cmd = x.textarea.Update(msg)
	x.textarea = &t
	return x, cmd
}

func (x xmlViewer) View() string {
	return x.textarea.View()
}

func (x xmlViewer) Width() int {
	return x.textarea.Width()
}

func (x xmlViewer) SetWidth(w int) {
	x.textarea.SetWidth(w)
}

func (x xmlViewer) SetHeight(h int) {
	x.textarea.SetHeight(h)
}

func (x xmlViewer) Focus() {
	x.textarea.Focus()
	x.textarea.Cursor.Focus()
	x.CursorTop()
}

func (x xmlViewer) CursorTop() {
	n := x.textarea.LineCount()
	for i := 0; i < n; i++ {
		x.textarea.CursorUp()
	}
	x.textarea.CursorStart()
}

func (x xmlViewer) Blur() {
	x.textarea.Blur()
}

func (x xmlViewer) Focused() bool {
	return x.textarea.Focused()
}

func (x xmlViewer) Highlight() {
	tag := false
	buffer := make([]rune, 0)
	out := ""
	for _, r := range x.textarea.Value() {
		if r == '<' {
			tag = true
		} else if r == '>' {
			tag = false
			out += lipgloss.NewStyle().Faint(true).Render("<" + string(buffer) + ">")
			buffer = buffer[:0]
		} else if tag {
			buffer = append(buffer, r)
		} else {
			out += string(r)
		}
	}
	x.textarea.SetValue(out)
}
