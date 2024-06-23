package app_old

import (
	"fmt"

	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/filelist"
	"github.com/JulianVidal/tagger/taglist"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Separate some of this stuff into files
// TODO: Current plan:
//			* Separate the file list and tag list into different files
//			* Add ways to create a tag
//			* Add ways to delete a tag
//			* Add a generic indexer so we can use locate, fzf, everything, etc...
//			* Way to show error codes from the engine

// NOTE: What I need
//   - A list of files
//   - Needs to be searchable through filename
//   - Needs to be searchable through tags
//   - Files should be taggable
//   - A list of tags
//   - Needs to be searchable through tag name
//   - Needs to be searchable throgh tags
//   - Tags should be taggable
var (
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	chosenTagsStyle = lipgloss.NewStyle().PaddingLeft(2)
)

type KeyMap struct {
	Switch key.Binding
	Print  key.Binding
	Edit   key.Binding
}

// TODO: Fix switch keybind showing up when adding tags to item
func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{m.KeyMap.Switch, m.KeyMap.Print}

	if m.tagFocus {
		kb = append(m.tagList.ShortHelp(), kb...)
	} else {
		kb = append(m.fileList.ShortHelp(), kb...)
	}

	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	first_row := []key.Binding{m.KeyMap.Switch, m.KeyMap.Print}
	var kb [][]key.Binding

	if m.tagFocus {
		first_row = append(m.tagList.FullHelp()[0], first_row...)
		kb = append([][]key.Binding{first_row}, m.tagList.FullHelp()[1:]...)
	} else {
		first_row = append(m.fileList.FullHelp()[0], first_row...)
		kb = append([][]key.Binding{first_row}, m.fileList.FullHelp()[1:]...)
	}

	return kb
}

var keys = KeyMap{
	Switch: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Filter Tags"),
	),
	Print: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Show engine"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit Tags"),
	),
}

func sliceEq[S []K, K comparable](a S, b S) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type Model struct {
	fileList    filelist.Model
	KeyMap      KeyMap
	help        help.Model
	tagList     taglist.Model
	editor      taglist.Model
	editorFocus bool
	tagFocus    bool
	engineFocus bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetTags(tags ...string) tea.Cmd {
	return tea.Batch(
		m.tagList.SetTags(tags...),
		m.editor.SetTags(tags...),
	)
}

func (m *Model) SetFiles(files []string) tea.Cmd {
	var fileItems []list.Item
	for _, file := range files {
		fileItems = append(fileItems, filelist.Item{Title: file})
	}
	return m.fileList.List.SetItems(fileItems)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.engineFocus {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.KeyMap.Print):
				m.engineFocus = !m.engineFocus
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Edit):
			if m.tagFocus {
				tag := m.tagList.List.SelectedItem().(taglist.Item)
				if m.editorFocus {
					tag.Tags = m.editor.ChosenTags()
					handler.SetTagParents(tag.Title, tag.Tags)
				} else {
					tag.Tags = handler.TagParents(tag.Title)
					m.editor.SetTags(handler.ValidParenTags(tag.Title)...)
					m.editor.SetChosen(tag.Tags...)
				}
				m.tagList.List.SetItem(m.tagList.List.Index(), tag)
			} else {
				item := m.fileList.List.SelectedItem().(filelist.Item)
				if m.editorFocus {
					item.Tags = m.editor.ChosenTags()
					handler.SetObjectTags(item.Title, item.Tags)
				} else {
					item.Tags = handler.ObjectTags(item.Title)
					m.editor.SetTags(handler.Tags()...)
					m.editor.SetChosen(item.Tags...)
				}
				m.fileList.List.SetItem(m.fileList.List.Index(), item)
			}
			m.editorFocus = !m.editorFocus
		case key.Matches(msg, m.KeyMap.Print):
			if m.fileList.List.FilterState() != list.Filtering &&
				m.tagList.List.FilterState() != list.Filtering {
				m.engineFocus = !m.engineFocus
			}
		case key.Matches(msg, m.KeyMap.Switch):
			if m.fileList.List.FilterState() != list.Filtering &&
				m.tagList.List.FilterState() != list.Filtering &&
				!m.editorFocus {
				m.tagFocus = !m.tagFocus
				if !m.tagFocus {
					chosenTags := m.tagList.ChosenTags()
					if len(chosenTags) != 0 {
						m.SetFiles(handler.QueryEngine(chosenTags))
						break
					}
					m.fileList.List.SetItems(m.fileList.Files)
				}
			}
		}
	}

	if m.editorFocus {
		m.editor, cmd = m.editor.Update(msg)
		m.help.ShowAll = m.editor.Help.ShowAll
	} else {
		if m.tagFocus {
			m.tagList, cmd = m.tagList.Update(msg)
			m.help.ShowAll = m.tagList.Help.ShowAll
		} else {
			m.fileList, cmd = m.fileList.Update(msg)
			m.help.ShowAll = m.fileList.Help.ShowAll
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.engineFocus {
		return handler.EngineString()
	}
	chosenTags := chosenTagsStyle.Render(
		fmt.Sprintf("Selected Tags: %s\n", m.tagList.ChosenTags()),
	)

	var focusedList string
	if m.tagFocus {
		focusedList = m.tagList.View()
	} else {
		focusedList = m.fileList.View()
	}
	if m.editorFocus {
		focusedList = lipgloss.JoinHorizontal(
			lipgloss.Left, focusedList, m.editor.View())
	}

	helpView := helpStyle.Render(m.help.View(m))
	return lipgloss.JoinVertical(lipgloss.Left, chosenTags, focusedList, helpView)
}

func New() Model {
	fl := filelist.New()
	tl := taglist.New()
	ed := taglist.New()
	return Model{fileList: fl, KeyMap: keys, help: help.New(), tagList: tl, editor: ed}
}
