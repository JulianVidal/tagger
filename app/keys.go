package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	FilePage key.Binding
	TagPage  key.Binding
	Print    key.Binding
	Quit     key.Binding
}

var keys = KeyMap{
	FilePage: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "Files"),
	),
	TagPage: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "Tags"),
	),
	Print: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Show engine"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
}
