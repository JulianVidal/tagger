package filepage

import (
	"github.com/JulianVidal/tagger/app/editor"
	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/app/taglist"
	"github.com/JulianVidal/tagger/internal/indexer"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	KeyMap    KeyMap
	title     string
	focus     Panel
	fileList  list.Model
	tagFilter taglist.Model
	editor    editor.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.fileList.SetWidth(msg.Width)
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
	case TagFilter:
		m.tagFilter, cmd = m.tagFilter.Update(msg)
	case FileList:
		m.fileList, cmd = m.fileList.Update(msg)
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
			m.tagFilter.KeyMap.Select) &&
			m.focus == TagFilter:

			tagged_files := handler.QueryEngine(m.tagFilter.ChosenTags())
			files := indexer.Query("")

			if len(tagged_files) != 0 {
				files = union(tagged_files, files)
			}

			cmd = m.SetFiles(files)
			cmds = append(cmds, cmd)

		case key.Matches(msg,
			m.fileList.KeyMap.CursorDown,
			m.fileList.KeyMap.CursorUp) &&
			m.focus == FileList:

			item := m.fileList.SelectedItem().(Item)
			m.editor.SetEditorObject(item.Title)

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
	tagFilterStyle := pageStyle
	fileListStyle := pageStyle
	editorStyle := pageStyle
	switch m.focus {
	case TagFilter:
		tagFilterStyle = pageFocusStyle
	case FileList:
		fileListStyle = pageFocusStyle
	case Editor:
		editorStyle = pageFocusStyle
	}
	return lipgloss.JoinHorizontal(lipgloss.Top,
		tagFilterStyle.Width(20).Render(m.tagFilter.View()),
		fileListStyle.Width(30).Render(m.fileList.View()),
		editorStyle.Width(30).Render(m.editor.View()),
	)
}

func New() Model {

	items := []list.Item{}
	for _, file := range indexer.Query("") {
		items = append(items, Item{Title: file})
	}

	l := list.New(items, itemDelegate{}, 20, 14)
	l.Title = "Files"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	ed := editor.New()

	fl := taglist.New()
	fl.List.Title = "Filter by tags"

	m := Model{
		KeyMap:    keys,
		title:     "Files",
		fileList:  l,
		editor:    ed,
		tagFilter: fl,
		focus:     FileList,
	}

	m.UpdateTags()
	return m
}
