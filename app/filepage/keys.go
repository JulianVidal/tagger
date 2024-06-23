package filepage

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	EnterEdit key.Binding
	ExitEdit  key.Binding
	Left      key.Binding
	Right     key.Binding
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
