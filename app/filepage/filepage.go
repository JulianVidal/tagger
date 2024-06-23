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

func union[S []K, K comparable](as S, bs S) S {
	a_set := make(map[K]struct{})
	for _, a := range as {
		a_set[a] = struct{}{}
	}
	cs := []K{}

	for _, b := range bs {
		if _, ok := a_set[b]; ok {
			cs = append(cs, b)
		}
	}
	return cs
}

type Model struct {
	KeyMap   KeyMap
	fileList list.Model
	title    string
	tagList  taglist.Model
	focus    Panel
	editor   editor.Model
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

		case key.Matches(msg, m.KeyMap.Left):
			m.switchPanel(m.focus.Prev())

		case key.Matches(msg, m.KeyMap.Right):
			m.switchPanel(m.focus.Next())
		}
	}

	switch m.focus {
	case TagFilter:
		m.tagList, cmd = m.tagList.Update(msg)
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
			m.tagList.KeyMap.Select) &&
			m.focus == TagFilter:

			tagged_files := handler.QueryEngine(m.tagList.ChosenTags())
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
	BorderStyle(lipgloss.RoundedBorder())

func (m Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		pageStyle.Width(20).Render(m.tagList.View()),
		pageStyle.Width(30).Render(m.fileList.View()),
		pageStyle.Width(30).Render(m.editor.View()),
	)
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

	tl := editor.New()
	item := l.SelectedItem().(Item)
	tl.SetEditorObject(item.Title)

	fl := taglist.New()
	fl.SetTags(handler.Tags()...)

	return Model{
		KeyMap:   keys,
		title:    "Files",
		fileList: l,
		editor:   tl,
		tagList:  fl,
		focus:    FileList,
	}
}
