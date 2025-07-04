package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	START_PAGE = iota
	FORM_PAGE
	SEARCH_PAGE
)

type page interface {
	update(msg tea.Msg) (tea.Cmd, int)
	view() string
}

type model struct {
	pages     [3]page
	pageIndex int
}

func initModel() model {
	return model{
		pages: [3]page{
			newStartPageModel(),
			newFormPageModel(),
			newSearchPageModel(),
		},
		pageIndex: START_PAGE,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	cmd, pageIndex := m.pages[m.pageIndex].update(msg)

	if m.pageIndex != pageIndex {
		if m.pageIndex == 2 && pageIndex == 1 {
			searchPageModel := m.pages[2].(*searchPageModel)
			fd := searchPageModel.searchResults[searchPageModel.recordIndex]

			m.pages[pageIndex] = formPageModelFromRecord(fd)
			m.pageIndex = pageIndex

			return m, cmd
		}

		var newPage page
		switch pageIndex {
		case START_PAGE:
			newPage = newStartPageModel()
		case FORM_PAGE:
			newPage = newFormPageModel()
		case SEARCH_PAGE:
			newPage = newSearchPageModel()
		}

		m.pages[pageIndex] = newPage
		m.pageIndex = pageIndex
	}

	return m, cmd
}

func (m model) View() string {
	page := m.pages[m.pageIndex]
	return page.view()
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
