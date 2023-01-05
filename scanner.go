package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// dummy component to fetch data
type scanner struct {
	directory string
	busy      bool
}

type ScanResultsMsg FeedbackResults
type ScanTriggerMsg string

func NewScanner(directory string) scanner {
	return scanner{directory: directory, busy: false}
}

func (s scanner) Init() tea.Cmd {
	return nil
}

func (s scanner) Update(msg tea.Msg) (scanner, tea.Cmd) {
	switch msg.(type) {
	case ScanTriggerMsg:
		return s, s.scan
	}
	return s, nil
}

func (s scanner) View() string {
	return ""
}

func (s scanner) rawScan() FeedbackResults {
	s.busy = true
	fr := search(s.directory)
	s.busy = false
	return fr
}

func (s scanner) scan() tea.Msg {
	return ScanResultsMsg(s.rawScan())
}
