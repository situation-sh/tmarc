package main

import "github.com/charmbracelet/lipgloss"

type colorTheme struct {
	primary lipgloss.Color
}

var themes = map[string]colorTheme{
	"default": {primary: lipgloss.Color("202")},
	"teal":    {primary: lipgloss.Color("#1D8489")},
	"pink":    {primary: lipgloss.Color("#E7388A")},
}

func ListThemes() []string {
	out := make([]string, 0)
	for t := range themes {
		out = append(out, t)
	}
	return out
}

func Theme() colorTheme {
	return themes[selectedTheme]
}
