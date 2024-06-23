package editor

import (
	"github.com/JulianVidal/tagger/app/taglist"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorItem interface {
	Parents() []string
	SetParent(string, bool)
	PossibleParents() []string
}

type Model struct {
	editorItem EditorItem
	tagList    taglist.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.tagList, cmd = m.tagList.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.tagList.KeyMap.Select):
			item := m.tagList.List.SelectedItem().(taglist.Item)
			m.editorItem.SetParent(item.Title, item.Selected)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.tagList.View()
}

func New() Model {
	tl := taglist.New()
	tl.List.Title = "Add tags"

	m := Model{tagList: tl}

	return m
}
