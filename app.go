package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Separate some of this stuff into files
// TODO: Current plan:
//			* Separate the file list and tag list into different files
//			* Add engine support
//			* Separate page for adding and removing tags

type KeyMap struct {
	Select key.Binding
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{m.KeyMap.Select}
	kb = append(m.list.ShortHelp(), kb...)
	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	first_row := []key.Binding{m.KeyMap.Select}
	first_row = append(m.list.FullHelp()[0], first_row...)
	kb := append([][]key.Binding{first_row}, m.list.FullHelp()[1:]...)
	return kb
}

var keys = KeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title string
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(strs ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(strs, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (item item) String() string {
	return fmt.Sprintf("%s", item.title)
}

type Model struct {
	list     list.Model
	KeyMap   KeyMap
	help     help.Model
	tagList  list.Model
	tagFocus bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Select):
			m.tagFocus = !m.tagFocus
			cmds = append(cmds, textinput.Blink)
		}
	}

	if m.tagFocus {
		m.tagList, cmd = m.tagList.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	m.help.ShowAll = m.list.Help.ShowAll
	help := helpStyle.Render(m.help.View(m))
	view := ""

	if m.tagFocus {
		view += m.tagList.View()
	} else {
		view += m.list.View()
	}

	view += help
	return view
}

func crawlDir(dir string) []string {
	var filepaths []string

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filepaths = append(filepaths, path)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	return filepaths
}

func main() {

	dir := os.Args[1]

	items := []list.Item{}
	for _, path := range crawlDir(dir) {
		items = append(items, item{
			title: path,
		})
	}

	const defaultWidth = 20
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowHelp(false)
	l.Title = "Search your files"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	tl := list.New([]list.Item{
		item{title: "bt"},
	}, itemDelegate{}, defaultWidth, listHeight)
	tl.SetShowHelp(false)
	tl.Title = "Search your tags"
	tl.SetShowStatusBar(true)
	tl.SetFilteringEnabled(true)
	tl.Styles.Title = titleStyle
	tl.Styles.PaginationStyle = paginationStyle
	tl.Styles.HelpStyle = helpStyle

	m := Model{list: l, KeyMap: keys, help: help.New(), tagList: tl}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
