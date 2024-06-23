package tagpage

import (
	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/app/taglist"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	KeyMap      KeyMap
	tagList     taglist.Model
	title       string
	editor      taglist.Model
	editorFocus bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.IsFiltering():

		case key.Matches(msg, m.KeyMap.EnterEdit) && !m.editorFocus:
			item, ok := m.tagList.List.SelectedItem().(taglist.Item)
			if ok {
				item.Tags = handler.TagParents(item.Title)
				m.editor.SetTags(handler.ValidParenTags(item.Title)...)
				m.editor.SetChosen(item.Tags...)
			}
			m.editorFocus = !m.editorFocus
			return m, nil

		case key.Matches(msg, m.KeyMap.ExitEdit) && m.editorFocus:
			item := m.tagList.List.SelectedItem().(taglist.Item)
			item.Tags = m.editor.ChosenTags()
			handler.SetTagParents(item.Title, item.Tags)
			m.editorFocus = !m.editorFocus
			return m, nil
		}

	}

	if m.editorFocus {
		m.editor, cmd = m.editor.Update(msg)
	} else {
		m.tagList.List, cmd = m.tagList.List.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var pageStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Width(50)

func (m Model) View() string {
	editor := ""
	if m.editorFocus {
		editor = pageStyle.Render(m.editor.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		pageStyle.Render(m.tagList.View()),
		editor)
}

func New() Model {
	tl := taglist.New()
	ed := taglist.New()

	return Model{KeyMap: keys, tagList: tl, title: "Tags", editor: ed}
}
