package tagpage

import tea "github.com/charmbracelet/bubbletea"

func (m Model) IsFiltering() bool {
	return m.tagList.IsFiltering()
}

func (m Model) Title() string {
	return m.title
}

func (m Model) ChosenTags() []string {
	return m.tagList.ChosenTags()
}

func (m *Model) SetTags(tags ...string) tea.Cmd {
	return m.tagList.SetTags(tags...)
}
