package mainmenu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/christianhturner/go-claude/tui/constants"
)

func selectMenuItem(selection constants.SessionState) tea.Cmd {
	return func() tea.Msg {
		return SessionStateMsg{SessionState: selection}
	}
}
