package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/junwei890/crawler/ui"
)

func main() {
	if _, err := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println(err)
	}
}
