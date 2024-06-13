package tagpage

import (
	"github.com/JulianVidal/tagger/app/taglist"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	tagList taglist.Model
	title   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.tagList, cmd = m.tagList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.tagList.View()
}

func New() Model {
	tl := taglist.New()

	return Model{tagList: tl, title: "Tags"}
}
