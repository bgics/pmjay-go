package starter

import (
	"fmt"

	"github.com/bgics/pmjay-go/internal/tui"
	"github.com/bgics/pmjay-go/internal/tui/view"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ExitError    error
	currentModel tea.Model
	sharedState  *tui.SharedState
}

func NewModel() *Model {
	s := tui.NewSharedState()
	return &Model{
		currentModel: view.NewStartPageModel(s),
		sharedState:  s,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.currentModel.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tui.FatalErrorMsg:
		m.ExitError = msg.Err
		return m, tea.Quit
	case tui.ErrorMsg:
		m.sharedState.Error = msg.Err
		return m, nil
	case tui.ChangePageMsg:
		if err := m.changePage(msg.To); err != nil {
			return m, tui.FatalErrorCmd(err)
		}

		cmd := m.currentModel.Init()

		return m, cmd
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	return m.currentModel.View()
}

func (m *Model) changePage(to tui.PageIndex) error {
	m.sharedState.Error = nil
	switch to {
	case tui.START_PAGE:
		m.currentModel = view.NewStartPageModel(m.sharedState)
		return nil
	case tui.SEARCH_PAGE:
		m.currentModel = view.NewSearchPageView(m.sharedState)
		return nil
	case tui.FORM_PAGE:
		m.currentModel = view.NewFormPageModel(m.sharedState)
		return nil
	}

	return fmt.Errorf("invalid page index %d", to)
}
