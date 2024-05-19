package taglist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Chosen key.Binding
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{m.KeyMap.Chosen}
	kb = append(m.List.ShortHelp(), kb...)
	return kb
}

func (m Model) FullHelp() [][]key.Binding {
	first_row := []key.Binding{m.KeyMap.Chosen}
	first_row = append(m.List.FullHelp()[0], first_row...)
	kb := append([][]key.Binding{first_row}, m.List.FullHelp()[1:]...)
	return kb
}

var keys = KeyMap{
	Chosen: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	chosenItemStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Item struct {
	Title    string
	Selected bool
	Tags     []string
}

func (i Item) FilterValue() string { return i.Title }
func (item Item) String() string {
	return fmt.Sprintf("%s", item.Title)
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	str := ""
	if i.Selected {
		str = fmt.Sprintf("[X] %s", i)
	} else {
		str = fmt.Sprintf("[ ] %s", i)
	}

	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(strs ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(strs, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (m Model) ChosenTags() []string {
	tags := []string{}

	for _, tag := range m.List.Items() {
		tag := tag.(Item)
		if tag.Selected {
			tags = append(tags, tag.Title)
		}
	}

	return tags
}

func (m *Model) SetTags(tags ...string) tea.Cmd {
	var items []list.Item
	for _, tag := range tags {
		items = append(items, Item{Title: tag})
	}
	return m.List.SetItems(items)
}

func (m *Model) SetChosen(tags ...string) {
	m.clearChosen()
	for i, item := range m.List.Items() {
		item := item.(Item)
		for _, tag := range tags {
			if item.Title == tag {
				item.Selected = true
				m.List.SetItem(i, item)
			}
		}
	}
}

func (m *Model) clearChosen() {
	for i, tag := range m.List.Items() {
		tag := tag.(Item)
		if tag.Selected {
			tag.Selected = false
			m.List.SetItem(i, tag)
		}
	}
}

type Model struct {
	List   list.Model
	KeyMap KeyMap
	Help   help.Model
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
		case key.Matches(msg, m.KeyMap.Chosen):
			tagItem := m.List.SelectedItem().(Item)
			tagItem.Selected = !tagItem.Selected

			cmd = m.List.SetItem(m.List.Index(), tagItem)
			cmds = append(cmds, cmd)
		}
	}

	m.List, cmd = m.List.Update(msg)
	m.Help.ShowAll = m.List.Help.ShowAll

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	view := ""
	view += m.List.View()
	return view
}

func New() Model {
	const defaultWidth = 20

	l := list.New([]list.Item{Item{Title: "test"}}, itemDelegate{}, defaultWidth, listHeight)

	l.SetShowHelp(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	l.Title = "Filter with tags"

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("100"))

	return Model{List: l, KeyMap: keys, Help: help.New()}
}
