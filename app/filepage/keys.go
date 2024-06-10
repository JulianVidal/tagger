package filepage

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Edit key.Binding
}

var keys = KeyMap{
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "Edit tags"),
	),
}
