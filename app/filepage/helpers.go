package filepage

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) IsFiltering() bool {
	return m.fileList.FilterState() == list.Filtering ||
		m.editor.IsFiltering() ||
		m.tagList.IsFiltering()
}

func (m Model) Title() string {
	return m.title
}

func (m *Model) SetFiles(files []string) tea.Cmd {
	var fileItems []list.Item
	for _, file := range files {
		fileItems = append(fileItems, Item{Title: file})
	}
	return m.fileList.SetItems(fileItems)
}

func mod(a, b int) int {
	return (a%b + b) % b
}
