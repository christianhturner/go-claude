package cliui

import (
	"fmt"

	"github.com/christianhturner/go-claude/chat"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
)

func ShowHistoricMessagePairs(displayAmount int, messagePairs []chat.MessagePair) {
	if displayAmount == 0 {
		fmt.Println("Nothing to print")
		return
	}
	maxPairs := len(messagePairs)
	if displayAmount > maxPairs {
		logger.Error(fmt.Sprintf("Display amount provided, %d, is greater than the number of message pairs, %d.\n", displayAmount, maxPairs))
	}
	fmt.Println("\nConversationHistory:")
	for i := 0; i < displayAmount && i < maxPairs; i++ {
		pair := messagePairs[i]
		fmt.Printf("\nUser: %s\n", pair.UserMessage.Content)
		fmt.Printf("\nClaude: %s\n", pair.AssistantMessage.Content)
	}
	fmt.Println("\n")
}

func ShowConvList(conv []db.Conversation) error {
	term := terminal.New()

	table := term.NewTable(30)

	idMaxWidth := 7
	table.AddColumn("ID", "ID", 5, &idMaxWidth, false, 0)
	table.AddColumn("Title", "Title", 40, nil, true, 0)
	timeMaxWidth := 19
	table.AddColumn("Created", "Created", 19, &timeMaxWidth, false, terminal.AlignCenter)
	table.AddColumn("Updated", "Updated", 19, &timeMaxWidth, false, terminal.AlignCenter)

	for _, c := range conv {
		table.AddRow(map[string]interface{}{
			"ID":      c.ID,
			"Title":   c.Title,
			"Created": c.CreatedAt.Format("2006-01-02 15:04:05"),
			"Updated": c.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	table.Render()

	return nil
}
