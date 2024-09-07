package conversationmenu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/christianhturner/go-claude/db"
	go_claude_list "github.com/christianhturner/go-claude/list"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/tui/constants"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	cmd      tea.Cmd
)

type BackMsg bool

type conversationItem struct {
	title       string
	description string
	// order       int
	// id          int64
}

func (i conversationItem) Title() string       { return i.title }
func (i conversationItem) Description() string { return i.description }
func (i conversationItem) FilterValue() string { return i.title }

type ConversationMenuModel struct {
	list list.Model
}

func (m ConversationMenuModel) Init() tea.Cmd {
	return nil
}

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(lipgloss.Color("170"))
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(lipgloss.Color("240"))
	return d
}

func InitialConversationMenu() tea.Model {
	menuItems := []list.Item{
		conversationItem{
			title:       "Title 1",
			description: "Some descriptiong",
		},
	}
	conversationMenu := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	conversationMenu.Title = "Test"
	return ConversationMenuModel{
		list: conversationMenu,
	}
	// items, err := newConversationList()
	// if err != nil {
	// 	return nil, func() tea.Msg { return errMsg{err} }
	// }
	// logger.Info(fmt.Sprintf("Loaded %d conversation items", len(items)))
	// delegate := newItemDelegate()
	// m := conversationMenuModel{list: list.New(items, delegate, 0, 0)}
	// m.list.Title = "Conversations"
	// m.list.SetShowStatusBar(false)
	// m.list.SetFilteringEnabled(false)
	// logger.Info(fmt.Sprintf("Initialized list with %d items", len(m.list.Items())))
	// logger.Info("Populating conversation tables")
	// return m, nil
}

func (m ConversationMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, constants.SelectMenuItem(constants.MainMenuView)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		// case updateConversationListMsg:
		// 	m.list = msg.Model
		// 	logger.Info(fmt.Sprintf("Updated list with %d items", len(m.list.Items())))
		// 	return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ConversationMenuModel) View() string {
	logger.Info(fmt.Sprintf("Rendering conversation menu with %d items", len(m.list.Items())))
	logger.Info(fmt.Sprintf("First Item: %v", m.list.Items()[0]))
	view := m.list.View()
	logger.Info(fmt.Sprintf("List View: \n%s", view))
	return docStyle.Render(m.list.View())
}

func (m *ConversationMenuModel) SetSize(width, height int) {
	logger.Info(fmt.Sprintf("Setting conversation menu size to %dx%d", width, height))
	m.list.SetSize(width, height)
}

func newConversationList() ([]list.Item, error) {
	conv, err := go_claude_list.GetAllConversations()
	if err != nil {
		logger.LogError(err, fmt.Sprintf("Error getting conversations from db: %v", err))
	}
	return conversationsToItems(conv), err
}

func conversationsToItems(conversations []db.Conversation) []list.Item {
	items := make([]list.Item, len(conversations))
	for i, conv := range conversations {
		items[i] = conversationItem{
			title:       conv.Title,
			description: fmt.Sprintf("Conversation number %d", i), // TODO: Message Preview for the description?
			// order:       i,
			// id:          conv.ID,
		}
	}
	return items
}
