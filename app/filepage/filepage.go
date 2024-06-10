package filepage

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	KeyMap   KeyMap
	FileList list.Model
	Title    string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.FileList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Edit):
		}

	}

	m.FileList, cmd = m.FileList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.FileList.View()
}

func New() Model {

	items := []list.Item{
		Item{Title: "di"},
		Item{Title: "temporary title"},
		Item{Title: "tempotitle #2"},
	}

	l := list.New(items, itemDelegate{}, 20, 14)
	l.Title = "Files"
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)

	return Model{KeyMap: keys, FileList: l, Title: "Files"}
}
