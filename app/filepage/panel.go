package filepage

type Panel int

const (
	TagFilter Panel = iota
	FileList
	Editor
	max
)

func (p Panel) Next() Panel {
	return Panel(mod(int(p)+1, int(max)))
}

func (p Panel) Prev() Panel {
	return Panel(mod(int(p)-1, int(max)))
}

func (m *Model) switchPanel(nextPanel Panel) {
	switch m.focus {
	case Editor:
		// m.editor.ResetCursor()
		// m.editor.ClearTags()
	}

	m.focus = nextPanel

	switch m.focus {
	case Editor:
		// item := m.fileList.SelectedItem().(Item)
		// m.editor.SetEditorObject(item.Title)
	}
}
