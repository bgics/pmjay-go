package main

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.NoColor{}).
			MaxWidth(50).
			MaxHeight(3)

	formFieldNameStyle = lipgloss.NewStyle().
				Width(13).
				MarginRight(1).
				AlignHorizontal(lipgloss.Right)

	dateIndicatorStyle = lipgloss.NewStyle().
				PaddingLeft(1)

	dateFieldStyle = dateIndicatorStyle.
			MarginTop(1)

	numDaysFieldStyle = dateFieldStyle

	printBtnStyle = formFieldNameStyle.
			MarginTop(2)

	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("202")).
			MarginTop(2).
			MarginLeft(2)
)
