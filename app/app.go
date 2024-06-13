package app

import (
	"fmt"

	"github.com/JulianVidal/tagger/app/filepage"
	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/app/tagpage"
	"github.com/JulianVidal/tagger/internal/indexer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	FilePage = iota
	TagPage
	c2
)

type Page interface {
	tea.Model
	IsFiltering() bool
	Title() string
}

type Model struct {
	KeyMap   KeyMap
	Focus    Page
	FilePage filepage.Model
	TagPage  tagpage.Model
}

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

			tagged_files := handler.QueryEngine(m.TagPage.ChosenTags())
			files := indexer.Query("")

			if len(tagged_files) != 0 {
				files = union(tagged_files, files)
			}

			cmd = m.FilePage.SetFiles(files)
			cmds = append(cmds, cmd)

			m.Focus = m.FilePage
		case key.Matches(msg, m.KeyMap.TagPage):
			m.Focus = m.TagPage
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
	for i, page := range []Page{m.FilePage, m.TagPage} {
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

	return Model{
		KeyMap:   keys,
		Focus:    tp,
		FilePage: fp,
		TagPage:  tp,
	}
}
