package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	startPageChoices = [2]string{"New Patient", "Search Records"}
)

type startPageModel struct {
	chosenIndex int
}

func newStartPageModel() *startPageModel {
	return &startPageModel{}
}

func (m *startPageModel) update(msg tea.Msg) (tea.Cmd, int) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch m.chosenIndex {
			case 0:
				return nil, FORM_PAGE
			case 1:
				return nil, SEARCH_PAGE
			}
		case "up", "shift+tab":
			m.chosenIndex--
		case "down", "tab":
			m.chosenIndex++
		}
	}

	if m.chosenIndex > len(startPageChoices)-1 {
		m.chosenIndex = 0
	} else if m.chosenIndex < 0 {
		m.chosenIndex = len(startPageChoices) - 1
	}
	return nil, START_PAGE
}

func (m *startPageModel) view() string {
	var output strings.Builder

	output.WriteString("\n\n")

	for i, choice := range startPageChoices {
		if i == m.chosenIndex {
			output.WriteString("  > " + choice + "\n\n")
		} else {
			output.WriteString("    " + choice + "\n\n")
		}
	}

	return output.String()
}
