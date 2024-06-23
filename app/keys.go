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
		key.WithKeys("3"),
		key.WithHelp("3", "Show engine"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}
