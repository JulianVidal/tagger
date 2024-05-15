package filelist

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JulianVidal/tagger/app/handler"
	"github.com/JulianVidal/tagger/taglist"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	EditTags key.Binding
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{m.KeyMap.EditTags}
	kb = append(m.List.ShortHelp(), kb...)
	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	first_row := []key.Binding{m.KeyMap.EditTags}
	first_row = append(m.List.FullHelp()[0], first_row...)
	kb := append([][]key.Binding{first_row}, m.List.FullHelp()[1:]...)
	return kb
}

var keys = KeyMap{
	EditTags: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit Tags"),
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

type Item struct {
	Title string
	Tags  []string
}

func (i Item) FilterValue() string { return i.Title }
func (item Item) String() string   { return fmt.Sprintf("%s", item.Title) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
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

type Model struct {
	List            list.Model
	KeyMap          KeyMap
	Help            help.Model
	Directory       string
	Files           []list.Item
	TagList         taglist.Model
	FocusingTagList bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.EditTags):
			item := m.List.SelectedItem().(Item)
			if m.FocusingTagList {
				item.Tags = m.TagList.GetChosenTags()
				m.List.SetItem(m.List.Index(), item)
				handler.SetObjectTags(item.Title, item.Tags)
			} else {
				item.Tags = handler.ObjectTags(item.Title)
				m.List.SetItem(m.List.Index(), item)
				m.TagList.SetChosen(item.Tags...)
			}
			m.FocusingTagList = !m.FocusingTagList
		}
	}

	if m.FocusingTagList {
		m.TagList, cmd = m.TagList.Update(msg)
		m.Help.ShowAll = m.TagList.Help.ShowAll
	} else {
		m.List, cmd = m.List.Update(msg)
		m.Help.ShowAll = m.List.Help.ShowAll
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	view := ""
	if m.FocusingTagList {
		view += lipgloss.JoinHorizontal(lipgloss.Left, m.List.View(), m.TagList.View())
	} else {
		view += m.List.View()
	}
	return view
}

func CrawlDir(dir string) []list.Item {
	var files []string

	err := filepath.Walk(dir,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, file)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	var items []list.Item
	for _, path := range files {
		items = append(items, Item{
			Title: path,
		})
	}

	return items
}

func NewItem(title string) Item {
	return Item{
		Title: title,
		Tags:  []string{"test 1", "test 2"},
	}
}

func New() Model {
	dir := os.Args[1]
	items := CrawlDir(dir)

	const defaultWidth = 20
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.SetShowHelp(false)
	l.Title = "Search your files"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("100"))

	tl := taglist.New()
	tl.List.Title = "Add tags to item"
	return Model{List: l, KeyMap: keys, Help: help.New(), Directory: dir, Files: items, TagList: tl}
}
