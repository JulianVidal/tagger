package filepage

type Panel int

const (
	TagFilter Panel = iota
	FileList
	Editor
	max
)

func (p Panel) Next() Panel {
	if int(p)+1 >= int(max) {
		return p
	}
	return Panel(int(p) + 1)
}

func (p Panel) Prev() Panel {
	if int(p)-1 < 0 {
		return p
	}
	return Panel(int(p) - 1)
}
