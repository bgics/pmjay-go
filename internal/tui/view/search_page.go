package view

import (
	"fmt"
	"strings"

	"github.com/bgics/pmjay-go/config"
	"github.com/bgics/pmjay-go/internal/tui"
	"github.com/bgics/pmjay-go/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchPageModel struct {
	searchInput   textinput.Model
	recordIndex   int
	searchResults []model.FormData
	sharedState   *tui.SharedState
}

func NewSearchPageView(sharedState *tui.SharedState) *SearchPageModel {
	s := &SearchPageModel{}
	s.searchInput = makeTextInput(true, config.NAME)
	s.sharedState = sharedState

	return s
}

func (m *SearchPageModel) Init() tea.Cmd {
	return textinput.Blink
}

// TODO: refactor update and views
func (m *SearchPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			m.recordIndex = cyclicAdjust(m.recordIndex-1, 0, max(len(m.searchResults)-1, 0))
		case "down", "tab":
			m.recordIndex = cyclicAdjust(m.recordIndex+1, 0, max(len(m.searchResults)-1, 0))
		case "enter":
			if len(m.searchResults) > 0 {
				m.sharedState.SelectedRecord = m.searchResults[m.recordIndex]
				m.sharedState.LastPageIndex = tui.SEARCH_PAGE
				return m, tui.ChangePageCmd(tui.FORM_PAGE)
			}

			return m, nil
		case "delete":
			return m.handleRemoveRecord()
		case "esc":
			m.sharedState.LastPageIndex = tui.SEARCH_PAGE
			return m, tui.ChangePageCmd(tui.START_PAGE)
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	searchValue := strings.TrimSpace(m.searchInput.Value())
	if len(searchValue) > 0 {
		results, err := m.sharedState.Store.GetRecordsByName(searchValue)
		if err != nil {
			return m, tui.ErrorCmd(err)
		} else {
			m.searchResults = results
		}
	} else {
		m.searchResults = nil
	}

	return m, cmd
}

// TODO: refactor the style uses in this function
func (m *SearchPageModel) View() string {
	var output strings.Builder

	output.WriteString("\n")

	searchInput := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"NAME ",
		tui.InputActiveBorderStyle.Render(m.searchInput.View()),
	)

	var searchResultViews []string
	for i, result := range m.searchResults {
		style := lipgloss.NewStyle().MarginTop(1)

		if i == m.recordIndex {
			searchResultViews = append(searchResultViews, style.Render("> "+result.Name))
		} else {
			searchResultViews = append(searchResultViews, style.Foreground(tui.InactiveColor).Render("  "+result.Name))
		}
	}

	searchResults := lipgloss.NewStyle().
		MarginLeft(5).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			searchResultViews...,
		),
	)

	var errMsg string
	if err := m.sharedState.Error; err != nil {
		errMsg = tui.ErrStyle.Render(fmt.Sprintf("[ERROR] %v", err))
	}

	output.WriteString(
		lipgloss.NewStyle().
			MarginLeft(3).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					searchInput,
					searchResults,
					errMsg,
				),
			),
	)

	return output.String()
}

func (m *SearchPageModel) handleRemoveRecord() (tea.Model, tea.Cmd) {
	recordName := m.searchResults[m.recordIndex].Name

	if len(m.searchResults) > 0 {
		if err := m.sharedState.Store.RemoveRecord(recordName); err != nil {
			return m, tui.ErrorCmd(err)
		}

		m.searchResults = append(m.searchResults[:m.recordIndex], m.searchResults[m.recordIndex+1:]...)

		m.recordIndex--
		if m.recordIndex < 0 {
			m.recordIndex = 0
		}
	}

	return m, nil
}
