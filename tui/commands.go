package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg                    struct{ error }
	updateConversationListMsg struct{ list.Model }
	renameConversationMsg     []list.Item
	switchViewMsg             struct{ view string }
)

func switchViewCmd(view string) tea.Cmd {
	return func() tea.Msg {
		return switchViewMsg{view}
	}
}

type convItem struct {
	title       string
	description string
	order       int
	id          int64
}

func (i convItem) Title() string       { return i.title }
func (i convItem) Description() string { return i.title }
func (i convItem) FilterValue() string { return i.title }

// func ConversationMenuCmd() tea.Cmd {
// 	return func() tea.Msg {
// 		model, _ := InitialConversationMenu()
// 		if m, ok := model.(conversationMenuModel); ok {
// 			logger.Info(fmt.Sprintf("Return conversation menu with %d items", len(m.list.Items())))
// 			return updateConversationListMsg{m.list}
// 		}
// 		return errMsg{fmt.Errorf("An Error occurred")}
// 	}
// }
