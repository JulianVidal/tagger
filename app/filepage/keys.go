package filepage

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	EnterEdit key.Binding
	ExitEdit  key.Binding
}

var keys = KeyMap{
	EnterEdit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Edit"),
	),
	ExitEdit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Exit Edit"),
	),
}
