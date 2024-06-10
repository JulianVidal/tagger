package tagpage

import (
	"github.com/JulianVidal/tagger/app/taglist"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	TagList taglist.Model
	Title   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.TagList, cmd = m.TagList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.TagList.View()
}

func New() Model {

	tl := taglist.New()

	return Model{TagList: tl, Title: "Tags"}
}
