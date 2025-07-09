package tui

import "github.com/charmbracelet/lipgloss"

var (
	InactiveColor = lipgloss.Color("240")
	ErrorColor    = lipgloss.Color("202")

	BorderStyle = lipgloss.NormalBorder()
)

var (
	InputActiveBorderStyle = lipgloss.NewStyle().
				Border(BorderStyle).
				BorderForeground(lipgloss.NoColor{}).
				MaxWidth(50).
				MaxHeight(3)

	InputInactiveBorderStyle = InputActiveBorderStyle.
					Foreground(InactiveColor).
					BorderForeground(InactiveColor)

	FieldNameActiveStyle = lipgloss.NewStyle().
				Width(14).
				MarginRight(1).
				AlignHorizontal(lipgloss.Right)

	FieldNameInactiveStyle = FieldNameActiveStyle.
				Foreground(InactiveColor)

	DateIndicatorStyle = lipgloss.NewStyle().
				PaddingLeft(1)

	DateFieldActiveStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				MarginTop(1)

	DateFieldInactiveStyle = DateFieldActiveStyle.
				Foreground(InactiveColor)

	NumDaysFieldActiveStyle = lipgloss.NewStyle().
				Margin(1).
				MarginTop(1)

	NumDaysFieldInactiveStyle = NumDaysFieldActiveStyle.
					Foreground(InactiveColor)

	BtnActiveStyle = lipgloss.NewStyle().
			Border(BorderStyle).
			Padding(0, 1).
			MarginLeft(1).
			MarginTop(1)

	BtnInactiveStyle = BtnActiveStyle.
				Foreground(InactiveColor).
				BorderForeground(InactiveColor)

	ErrStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			MarginTop(2).
			MarginLeft(2)
)
