package tagpage

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left   key.Binding
	Right  key.Binding
	Create key.Binding
	Delete key.Binding
	Enter  key.Binding
	Escape key.Binding
}

var keys = KeyMap{
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Cancel"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Confirm"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Create"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete"),
	),
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Right"),
	),
}
