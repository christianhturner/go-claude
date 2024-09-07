package constants

import tea "github.com/charmbracelet/bubbletea"

type SessionState int

const (
	ConversationMenuView SessionState = iota
	MainMenuView
)

func SelectMenuItem(selection SessionState) tea.Cmd {
	return func() tea.Msg {
		return SessionStateMsg{SessionState: selection}
	}
}

type SessionStateMsg struct {
	SessionState SessionState
}
