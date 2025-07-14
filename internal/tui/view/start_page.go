package view

import (
	"fmt"

	"github.com/bgics/pmjay-go/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			m.choiceIndex = cyclicAdjust(m.choiceIndex-1, 0, len(choices)-1)
		case "down", "tab":
			m.choiceIndex = cyclicAdjust(m.choiceIndex+1, 0, len(choices)-1)
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

func (m *StartPageModel) View() string {
	rows := make([]string, len(choices)+1)
	for i, choice := range choices {
		style := lipgloss.NewStyle().MarginTop(2)
		if i == m.choiceIndex {
			rows[i] = style.Render("> " + choice)
		} else {
			rows[i] = style.Render("  " + choice)
		}
	}

	if err := m.sharedState.Error; err != nil {
		rows[len(choices)] = tui.ErrStyle.Render(fmt.Sprintf("[ERROR] %v", err))
	}

	return lipgloss.NewStyle().
		Margin(0, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			rows...,
		))

}
