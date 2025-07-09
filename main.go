package main

import (
	"fmt"
	"os"

	"github.com/bgics/pmjay-go/internal/tui/starter"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(starter.NewModel(), tea.WithAltScreen())
	exitModel, err := p.Run()
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
		os.Exit(1)
	}

	typedExitModel, ok := exitModel.(*starter.Model)
	if !ok {
		fmt.Println("failed to assert exit model type")
		os.Exit(1)
	}

	if err := typedExitModel.ExitError; err != nil {
		fmt.Printf("model exited with error: %v\n", err)
		os.Exit(1)
	}
}
