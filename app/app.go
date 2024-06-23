package app

import (
	"fmt"

	"github.com/JulianVidal/tagger/app/enginepage"
	"github.com/JulianVidal/tagger/app/filepage"
	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/app/tagpage"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Add filtering feature again :/
// TODO: Colour outline or remove colour when switching panels
// TODO: Find a way to keep the tags in order
// TODO: Separate some of this stuff into files
// TODO: Current plan:
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

type Page interface {
	tea.Model
	IsFiltering() bool
	Title() string
}

type Model struct {
	KeyMap     KeyMap
	Focus      Page
	FilePage   filepage.Model
	TagPage    tagpage.Model
	EnginePage enginepage.Model
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
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case m.Focus.IsFiltering():
		case key.Matches(msg, m.KeyMap.FilePage):
			m.Focus = m.FilePage
		case key.Matches(msg, m.KeyMap.TagPage):
			m.Focus = m.TagPage
		case key.Matches(msg, m.KeyMap.Print):
			m.Focus = m.EnginePage
		}

	}

	updatedPage, cmd := m.Focus.Update(msg)
	m.Focus = updatedPage.(Page)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var titleStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder())

var titleFocusStyle = titleStyle.Copy().
	Bold(true).
	BorderForeground(lipgloss.Color("63")).
	Foreground(lipgloss.Color("#7D56F4"))

var pageStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	Width(50)

func (m Model) View() string {
	titles := []string{}
	for i, page := range []Page{m.FilePage, m.TagPage, m.EnginePage} {
		title := fmt.Sprintf("%d. %s", i+1, page.Title())
		if m.Focus.Title() == page.Title() {
			titles = append(titles, titleFocusStyle.Render(title))
		} else {
			titles = append(titles, titleStyle.Render(title))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, titles...),
		m.Focus.View(),
	)
}

func New() Model {
	fp := filepage.New()
	tp := tagpage.New()
	tp.SetTags(handler.Tags()...)
	ep := enginepage.Model{}

	return Model{
		KeyMap:     keys,
		Focus:      fp,
		FilePage:   fp,
		TagPage:    tp,
		EnginePage: ep,
	}
}
