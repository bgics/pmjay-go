package tui

import "github.com/charmbracelet/lipgloss"

var (
	InactiveColor            = lipgloss.Color("240")
	ErrorColor               = lipgloss.Color("202")
	DatePickerHighlightColor = lipgloss.Color("208")

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

	DateFieldActiveStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				MarginTop(1)

	DateFieldInactiveStyle = DateFieldActiveStyle.
				Foreground(InactiveColor)

	SimpleFieldActiveStyle = lipgloss.NewStyle().
				Margin(1).
				MarginTop(1)

	SimpleFieldInactiveStyle = SimpleFieldActiveStyle.
					Foreground(InactiveColor)

	BtnActiveStyle = lipgloss.NewStyle().
			Border(BorderStyle).
			Padding(0, 1).
			MarginLeft(1).
			MarginTop(1)

	BtnInactiveStyle = BtnActiveStyle.
				Foreground(InactiveColor).
				BorderForeground(InactiveColor)

	DatePickerStyle = lipgloss.NewStyle().
			Border(BorderStyle).
			BorderForeground(lipgloss.NoColor{}).
			Padding(0, 2).
			MarginLeft(5)

	ErrStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			MarginTop(2).
			MarginLeft(2)
)
