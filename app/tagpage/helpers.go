package tagpage

import (
	"github.com/JulianVidal/tagger/internal/engine"
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

func (m *Model) UpdateTags() {
	items := []list.Item{}
	for _, tag := range engine.Tags() {
		items = append(items, Item{Title: tag})
	}
	m.tagList.SetItems(items)

	item := m.tagList.SelectedItem().(Item)
	m.editor.SetEditorTag(item.Title)
}
