package tagpage

import (
	"github.com/JulianVidal/tagger/app/editor"
	"github.com/JulianVidal/tagger/app/handler"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	KeyMap  KeyMap
	title   string
	focus   Panel
	tagList list.Model
	editor  editor.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.tagList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case m.IsFiltering():

		case key.Matches(msg, m.KeyMap.Left):
			m.focus = m.focus.Prev()

		case key.Matches(msg, m.KeyMap.Right):
			m.focus = m.focus.Next()
		}
	}

	switch m.focus {
	case TagList:
		m.tagList, cmd = m.tagList.Update(msg)
	case Editor:
		m.editor, cmd = m.editor.Update(msg)
	default:
		panic("Wrong Focus Number")
	}
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg,
			m.tagList.KeyMap.CursorDown,
			m.tagList.KeyMap.CursorUp) &&
			m.focus == TagList:

			item := m.tagList.SelectedItem().(Item)
			m.editor.SetEditorTag(item.Title)

		}
	}

	return m, tea.Batch(cmds...)
}

var pageStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Width(50)

var pageFocusStyle = pageStyle.Copy().
	BorderForeground(lipgloss.Color("63"))

func (m Model) View() string {
	tagListStyle := pageStyle
	editorStyle := pageStyle
	switch m.focus {
	case TagList:
		tagListStyle = pageFocusStyle
	case Editor:
		editorStyle = pageFocusStyle
	}
	return lipgloss.JoinHorizontal(lipgloss.Top,
		tagListStyle.Width(30).Render(m.tagList.View()),
		editorStyle.Width(30).Render(m.editor.View()),
	)

}

func New() Model {

	items := []list.Item{}
	for _, tag := range handler.Tags() {
		items = append(items, Item{Title: tag})
	}

	tl := list.New(items, itemDelegate{}, 20, 14)
	tl.Title = "Tags"
	tl.SetShowHelp(false)
	tl.SetShowStatusBar(false)
	tl.DisableQuitKeybindings()

	ed := editor.New()
	item := tl.SelectedItem().(Item)
	ed.SetEditorTag(item.Title)

	return Model{
		KeyMap:  keys,
		title:   "Tags",
		tagList: tl,
		editor:  ed,
		focus:   TagList,
	}
}
