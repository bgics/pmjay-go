package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchPageModel struct {
	searchInput   textinput.Model
	store         store
	pageError     error
	searchResults []formData
	recordIndex   int
}

func newSearchPageModel() *searchPageModel {
	s := &searchPageModel{}

	s.searchInput = makeSearchTextInput()
	s.pageError = s.store.loadRecords()

	return s
}

func makeSearchTextInput() textinput.Model {
	t := textinput.New()

	t.Cursor.SetMode(cursor.CursorStatic)

	t.Prompt = " "
	t.Width = 40
	t.Focus()

	return t
}

func (m *searchPageModel) update(msg tea.Msg) (tea.Cmd, int) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, START_PAGE
		case "enter":
			if len(m.searchResults) > 0 {
				return nil, FORM_PAGE
			}
		case "shift+tab":
			m.recordIndex--
		case "tab":
			m.recordIndex++
		}
	}

	if len(m.searchResults) > 0 {
		if m.recordIndex > len(m.searchResults)-1 {
			m.recordIndex = 0
		} else if m.recordIndex < 0 {
			m.recordIndex = len(m.searchResults) - 1
		}
	} else {
		m.recordIndex = 0
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	searchValue := strings.TrimSpace(m.searchInput.Value())

	if len(searchValue) > 0 {
		m.searchResults = m.store.getRecordsByName(searchValue)
	} else {
		m.searchResults = nil
	}

	return cmd, SEARCH_PAGE
}

func (m *searchPageModel) view() string {
	var output strings.Builder

	output.WriteString("\n")

	output.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, "  NAME ", borderStyle.Render(m.searchInput.View())))

	output.WriteString("\n\n")

	for i, result := range m.searchResults {
		if i == m.recordIndex {
			output.WriteString("       > " + result.name + "\n\n")
		} else {
			output.WriteString("         " + result.name + "\n\n")
		}
	}

	return output.String()
}
