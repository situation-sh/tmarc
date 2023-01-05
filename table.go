package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const maxColWidth = 22

func tableStyle() table.Styles {
	return table.Styles{
		Selected: lipgloss.NewStyle().Bold(true).Background(Theme().primary).Foreground(lipgloss.Color("#FFFFFF")),
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:     lipgloss.NewStyle().Padding(0, 1),
	}
}

func toTable(results FeedbackResults) ([]table.Column, []table.Row) {
	r := results[0]
	cols := r.Columns(0)

	rows := make([]table.Row, 0)
	columns := make([]table.Column, len(cols))

	widths := make([]int, len(cols))
	for _, r := range results {
		r := r.ToRow(0)
		for k, s := range r {
			ls := len(s)
			if ls > maxColWidth {
				r[k] = s[:(maxColWidth-1)] + "…"
				widths[k] = maxColWidth
			} else if ls > widths[k] {
				widths[k] = ls
			}

		}
		rows = append(rows, r)
	}

	w := 4
	for k, c := range cols {
		title := strings.ToUpper(c)
		if len(title) > widths[k] {
			widths[k] = len(title)
		}

		columns[k] = table.Column{
			Title: title, Width: widths[k],
		}

		w += widths[k] + 1
	}
	// fmt.Println(w)
	return columns, rows
}

func NewTable(results FeedbackResults) table.Model {
	columns, rows := toTable(results)
	return table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithWidth(80),
		table.WithHeight(10),
		table.WithStyles(tableStyle()),
	)
}
