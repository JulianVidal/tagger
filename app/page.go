package app

type Page int

const (
	FilePage Page = iota
	TagPage
	EnginePage
	max
)

func (p Page) Next() Page {
	if int(p)+1 >= int(max) {
		return p
	}
	return Page(int(p) + 1)
}

func (p Page) Prev() Page {
	if int(p)-1 < 0 {
		return p
	}
	return Page(int(p) - 1)
}

func (p Page) String() string {
	return [...]string{"Files", "Tags", "Engine"}[p]
}
