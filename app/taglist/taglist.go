package taglist

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	KeyMap KeyMap
	List   list.Model
	Title  string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Edit):
		case key.Matches(msg, m.KeyMap.Select):
			tagItem := m.List.SelectedItem().(Item)
			tagItem.Selected = !tagItem.Selected

			cmd = m.List.SetItem(m.List.Index(), tagItem)
			cmds = append(cmds, cmd)
		}

	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.List.View()
}

func New() Model {
	l := list.New([]list.Item{}, itemDelegate{}, 20, 14)
	l.Title = "Tags"
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	return Model{KeyMap: keys, List: l, Title: "Tags"}
}
