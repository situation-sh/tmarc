package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	flag.StringVar(&directory, "d", ".", "directory to scan")
	flag.StringVar(&selectedTheme, "t", "default", fmt.Sprintf("color theme (%s)", strings.Join(ListThemes(), ", ")))
	flag.BoolVar(&highlightXML, "p", false, "enable xml syntax highlighting (experimental)")
	flag.Parse()

	m := NewModel(directory)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
