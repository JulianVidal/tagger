package tagpage

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) IsFiltering() bool {
	return m.tagList.FilterState() == list.Filtering ||
		m.editor.IsFiltering()
}

func (m Model) Title() string {
	return m.title
}

func (m *Model) SetTags(tags ...string) tea.Cmd {
	var tagItems []list.Item
	for _, tag := range tags {
		tagItems = append(tagItems, Item{Title: tag})
	}
	return m.tagList.SetItems(tagItems)
}
