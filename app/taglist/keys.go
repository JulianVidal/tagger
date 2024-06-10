package taglist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Edit   key.Binding
	Select key.Binding
}

var keys = KeyMap{
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit tags"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select"),
	),
}
