package taglist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("#7D56F4"))
)

type Item struct {
	Title    string
	Selected bool
	Tags     []string
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
	str := ""
	if i.Selected {
		str = fmt.Sprintf("[✗] %s", i)
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
