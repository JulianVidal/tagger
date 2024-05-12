package main

import (
	"fmt"
	"os"

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
//			* Separate page for adding and removing tags

var (
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	chosenTagsStyle = lipgloss.NewStyle().PaddingLeft(2)
)

type KeyMap struct {
	Switch key.Binding
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{m.KeyMap.Switch}

	if m.tagFocus {
		kb = append(m.tagList.ShortHelp(), kb...)
	} else {
		kb = append(m.fileList.ShortHelp(), kb...)
	}

	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	first_row := []key.Binding{m.KeyMap.Switch}
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
}

type Model struct {
	fileList filelist.Model
	KeyMap   KeyMap
	help     help.Model
	tagList  taglist.Model
	tagFocus bool
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
		case key.Matches(msg, m.KeyMap.Switch):
			if m.fileList.List.FilterState() != list.Filtering &&
				m.tagList.List.FilterState() != list.Filtering {
				m.tagFocus = !m.tagFocus
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

func main() {

	fl := filelist.New()
	tl := taglist.New()
	m := Model{fileList: fl, KeyMap: keys, help: help.New(), tagList: tl}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
