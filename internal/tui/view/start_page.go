package view

import (
	"strings"

	"github.com/bgics/pmjay-go/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	choices = []string{"New Patient", "Search Records"}
)

type StartPageModel struct {
	choiceIndex int
	sharedState *tui.SharedState
}

func NewStartPageModel(sharedState *tui.SharedState) *StartPageModel {
	return &StartPageModel{
		choiceIndex: 0,
		sharedState: sharedState,
	}
}

func (m *StartPageModel) Init() tea.Cmd {
	return nil
}

func (m *StartPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			m.decrementChoiceIndex()
		case "down", "tab":
			m.incrementChoiceIndex()
		case "enter":
			switch m.choiceIndex {
			case 0:
				m.sharedState.LastPageIndex = tui.START_PAGE
				return m, tui.ChangePageCmd(tui.FORM_PAGE)
			case 1:
				m.sharedState.LastPageIndex = tui.START_PAGE
				return m, tui.ChangePageCmd(tui.SEARCH_PAGE)
			}
		}
	}

	return m, nil
}

func (m *StartPageModel) incrementChoiceIndex() {
	m.choiceIndex++

	if m.choiceIndex > len(choices)-1 {
		m.choiceIndex = 0
	}
}

func (m *StartPageModel) decrementChoiceIndex() {
	m.choiceIndex--

	if m.choiceIndex < 0 {
		m.choiceIndex = len(choices) - 1
	}
}

func (m *StartPageModel) View() string {
	var output strings.Builder

	output.WriteString("\n\n")

	for i, choice := range choices {
		if i == m.choiceIndex {
			output.WriteString("  > " + choice + "\n\n")
		} else {
			output.WriteString("    " + choice + "\n\n")
		}
	}

	return output.String()
}
