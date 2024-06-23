package enginepage

import (
	"github.com/JulianVidal/tagger/internal/engine"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return engine.String()
}

func (m Model) IsFiltering() bool {
	return false
}

func (m Model) Title() string {
	return "Engine"
}
