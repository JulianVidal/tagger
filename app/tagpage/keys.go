package tagpage

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Left  key.Binding
	Right key.Binding
}

var keys = KeyMap{
	Left: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Right"),
	),
}
