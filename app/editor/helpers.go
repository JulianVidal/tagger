package editor

func (m *Model) IsFiltering() bool {
	return m.tagList.IsFiltering()
}

func (m *Model) ResetCursor() {
	m.tagList.List.Select(0)
}

func (m *Model) SetEditorTag(item string) {
	m.editorItem = NewEditorTag(item)
	m.updateTagList()
}

func (m *Model) SetEditorObject(item string) {
	m.editorItem = NewEditorObject(item)
	m.updateTagList()
}

func (m *Model) updateTagList() {
	m.tagList.SetTags(m.editorItem.PossibleParents()...)
	m.tagList.SetChosen(m.editorItem.Parents()...)
}

func (m *Model) ClearTags() {
	m.tagList.ClearTags()
}
