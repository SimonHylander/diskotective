package cmd

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	Enter      key.Binding
	Delete     key.Binding
	Yes        key.Binding
	No         key.Binding
	Backspace  key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "open directory"),
		),
		Backspace: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "go back"),
		),
		Delete: key.NewBinding(
			key.WithKeys("delete", "d"),
			key.WithHelp("delete", "delete item"),
		),
		Yes: key.NewBinding(
			key.WithKeys("y"),
		),
		No: key.NewBinding(
			key.WithKeys("n"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown, k.Enter, k.Delete},
		{k.Help, k.Quit},
	}
}
