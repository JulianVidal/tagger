package app

import (
	"fmt"

	"github.com/JulianVidal/tagger/app/enginepage"
	"github.com/JulianVidal/tagger/app/filepage"
	"github.com/JulianVidal/tagger/app/tagpage"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Find a way to keep the tags in order
// TODO: Delete hanler?
// TODO: Separate some of this stuff into files
// TODO: Current plan:
//			* Add ways to create a tag <-
//				* Will need to update, filter, taglists...
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

type Model struct {
	KeyMap     KeyMap
	FilePage   filepage.Model
	TagPage    tagpage.Model
	EnginePage enginepage.Model
	Focus      Page
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

		case m.FilePage.IsFiltering() && m.Focus == FilePage:
		case m.TagPage.IsFiltering() && m.Focus == TagPage:
		case m.EnginePage.IsFiltering() && m.Focus == EnginePage:

		case key.Matches(msg, m.KeyMap.FilePage):
			m.FilePage.UpdateTags()
			m.Focus = FilePage
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.KeyMap.TagPage):
			m.TagPage.UpdateTags()
			m.Focus = TagPage
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.KeyMap.Print):
			m.Focus = EnginePage
			return m, tea.Batch(cmds...)
		}

	}

	switch m.Focus {
	case FilePage:
		m.FilePage, cmd = m.FilePage.Update(msg)
	case TagPage:
		m.TagPage, cmd = m.TagPage.Update(msg)
	case EnginePage:
		m.EnginePage, cmd = m.EnginePage.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

var titleStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder())

var titleFocusStyle = titleStyle.Copy().
	Bold(true).
	BorderForeground(lipgloss.Color("63")).
	Foreground(lipgloss.Color("#7D56F4"))

func (m Model) View() string {
	titles := []string{}
	focus := ""
	switch m.Focus {
	case FilePage:
		focus = m.FilePage.View()
	case TagPage:
		focus = m.TagPage.View()
	case EnginePage:
		focus = m.EnginePage.View()
	}

	for p := FilePage; p < max; p++ {
		title := fmt.Sprintf("%d. %s", p+1, p.String())
		style := titleStyle
		if m.Focus == p {
			style = titleFocusStyle
		}
		titles = append(titles, style.Render(title))
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, titles...),
		focus,
	)
}

func New() Model {
	fp := filepage.New()
	tp := tagpage.New()
	ep := enginepage.Model{}

	return Model{
		KeyMap:     keys,
		Focus:      FilePage,
		FilePage:   fp,
		TagPage:    tp,
		EnginePage: ep,
	}
}
