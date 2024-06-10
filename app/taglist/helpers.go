package taglist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
