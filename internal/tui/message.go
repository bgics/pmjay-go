package tui

import tea "github.com/charmbracelet/bubbletea"

const (
	START_PAGE = iota
	SEARCH_PAGE
	FORM_PAGE
)

type PageIndex int

type ChangePageMsg struct {
	To PageIndex
}

type FatalErrorMsg struct {
	Err error
}

type ErrorMsg struct {
	Err error
}

func ChangePageCmd(to PageIndex) tea.Cmd {
	return func() tea.Msg {
		return ChangePageMsg{
			To: to,
		}
	}
}

func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg{
			Err: err,
		}
	}
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{
			Err: err,
		}
	}
}
