package filepage

import (
	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/app/taglist"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	KeyMap      KeyMap
	fileList    list.Model
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
	case tea.WindowSizeMsg:
		m.fileList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case m.IsFiltering():

		case key.Matches(msg, m.KeyMap.EnterEdit) && !m.editorFocus:
			item := m.fileList.SelectedItem().(Item)
			item.Tags = handler.ObjectTags(item.Title)
			m.editor.SetTags(handler.Tags()...)
			m.editor.SetChosen(item.Tags...)
			m.editorFocus = !m.editorFocus
			return m, nil

		case key.Matches(msg, m.KeyMap.ExitEdit) && m.editorFocus:
			item := m.fileList.SelectedItem().(Item)
			item.Tags = m.editor.ChosenTags()
			handler.SetObjectTags(item.Title, item.Tags)
			m.editorFocus = !m.editorFocus
			return m, nil
		}

	}

	if m.editorFocus {
		m.editor, cmd = m.editor.Update(msg)
	} else {
		m.fileList, cmd = m.fileList.Update(msg)
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
		editor = m.editor.View()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		pageStyle.Render(m.fileList.View()),
		pageStyle.Render(editor))
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
	l.DisableQuitKeybindings()

	tl := taglist.New()
	tl.SetTags(handler.Tags()...)

	return Model{KeyMap: keys, fileList: l, title: "Files", editor: tl}
}
