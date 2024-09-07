package go_claude_list

import (
	"fmt"

	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
)

func GetAllConversations() ([]db.Conversation, error) {
	return db.ListConversations()
}

func ShowMessageList(conversationId int64) error {
	messages, err := db.GetMessages(conversationId)
	if err != nil {
		logger.PanicError(err, "Error getting messages from DB")
	}
	conversations, err := db.GetConversation(conversationId)
	title := conversations.Title

	term := terminal.New()

	table := term.NewTable(30)

	idMaxWidth := 7
	table.AddColumn("ID", "ID", 5, &idMaxWidth, false, 0)
	table.AddColumn("Role", "Role", 5, nil, false, 0)
	table.AddColumn("Content", "Content", 40, nil, true, 0)
	table.AddColumn("Created", "Created", 19, nil, false, 0)

	for _, m := range messages {
		table.AddRow(map[string]interface{}{
			"ID":      m.ID,
			"Role":    m.Role,
			"Content": m.Content,
			"Created": m.CreatedAt,
		})
	}
	fmt.Printf("\nConversation: %s\n", title)
	table.Render()

	return nil
}
