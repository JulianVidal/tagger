package app

import (
	"fmt"

	"github.com/JulianVidal/tagger/app/filepage"
	"github.com/JulianVidal/tagger/app/tagpage"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	FilePage = iota
	TagPage
	c2
)

type Page struct {
	model tea.Model
	title string
}

type Model struct {
	KeyMap    KeyMap
	Pages     []Page
	PageIndex int
}

func (m Model) Init() tea.Cmd {
	return nil
}

type PageModel interface {
	tea.Model
	Name() string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.FilePage):
			m.PageIndex = 0
		case key.Matches(msg, m.KeyMap.TagPage):
			m.PageIndex = 1
		}

	}

	m.Pages[m.PageIndex].model, cmd = m.Pages[m.PageIndex].model.Update(msg)
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
	for i, page := range m.Pages {
		title := fmt.Sprintf("%d. %s", i+1, page.title)
		if i == m.PageIndex {
			titles = append(titles, titleFocusStyle.Render(title))
		} else {
			titles = append(titles, titleStyle.Render(title))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, titles...),
		pageStyle.Render(m.Pages[m.PageIndex].model.View()),
	)
}

func New() Model {
	fp := filepage.New()
	tp := tagpage.New()
	pages := []Page{
		{
			model: fp,
			title: fp.Title,
		},
		{
			model: tp,
			title: tp.Title,
		},
	}
	return Model{
		KeyMap: keys,
		Pages:  pages,
	}
}
