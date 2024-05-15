package app

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
//			* Add engine support
//			* For each item create a list of tags it can select from

var (
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	chosenTagsStyle = lipgloss.NewStyle().PaddingLeft(2)
)

type KeyMap struct {
	Switch key.Binding
	Print  key.Binding
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
	tagFocus    bool
	engineFocus bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetTags(tags ...string) tea.Cmd {
	return tea.Batch(
		m.tagList.SetTags(tags...),
		m.fileList.TagList.SetTags(tags...),
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
		case key.Matches(msg, m.KeyMap.Print):
			m.engineFocus = !m.engineFocus
		case key.Matches(msg, m.KeyMap.Switch):
			if m.fileList.List.FilterState() != list.Filtering &&
				m.tagList.List.FilterState() != list.Filtering &&
				!m.fileList.FocusingTagList {
				m.tagFocus = !m.tagFocus
				if !m.tagFocus {
					chosenTags := m.tagList.GetChosenTags()
					if len(chosenTags) != 0 {
						m.SetFiles(handler.QueryEngine(chosenTags))
						break
					}
					m.fileList.List.SetItems(m.fileList.Files)
				}
			}
		}
	}

	if m.tagFocus {
		m.tagList, cmd = m.tagList.Update(msg)
		m.help.ShowAll = m.tagList.Help.ShowAll
	} else {
		m.fileList, cmd = m.fileList.Update(msg)
		m.help.ShowAll = m.fileList.Help.ShowAll
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.engineFocus {
		return handler.EngineString()
	}
	chosenTags := chosenTagsStyle.Render(
		fmt.Sprintf("Selected Tags: %s\n", m.tagList.GetChosenTags()),
	)

	var focusedList string
	if m.tagFocus {
		focusedList = m.tagList.View()
	} else {
		focusedList = m.fileList.View()
	}

	helpView := helpStyle.Render(m.help.View(m))
	return lipgloss.JoinVertical(lipgloss.Left, chosenTags, focusedList, helpView)
}

func New() Model {
	fl := filelist.New()
	tl := taglist.New()
	return Model{fileList: fl, KeyMap: keys, help: help.New(), tagList: tl}
}
