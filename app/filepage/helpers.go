package filepage

import (
	"github.com/JulianVidal/tagger/internal/engine"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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

func (m Model) IsFiltering() bool {
	return m.fileList.FilterState() == list.Filtering ||
		m.editor.IsFiltering() ||
		m.tagFilter.IsFiltering()
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

func (m *Model) UpdateTags() {
	m.tagFilter.SetTags(engine.Tags()...)
	item := m.fileList.SelectedItem().(Item)
	m.editor.SetEditorObject(item.Title)
}

func (m Model) getFiles() []string {
	var tags []*engine.Tag
	for _, tagName := range m.tagFilter.ChosenTags() {
		tag, exists := engine.FindTag(tagName)
		if !exists {
			panic("Couldn't find tag in engine")
		}
		tags = append(tags, tag)
	}
	objects, err := engine.Query(tags...)
	if err != nil {
		panic(err)
	}

	var tagged_files []string
	for _, object := range objects {
		tagged_files = append(tagged_files, object.Name())
	}

	return tagged_files
}
