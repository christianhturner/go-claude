package delete

import (
	"fmt"

	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
)

func PromptForConversationId() int64 {
	conversations, err := db.ListConversations()
	options := make(map[interface{}]string)
	for _, convOptions := range conversations {
		options[convOptions.ID] = convOptions.Title
	}
	userSelect := terminal.New().PromptOptionsSelect(options)
	id, ok := userSelect.ID.(int64)
	if !ok {
		logger.PanicError(err, "Invalid ID type; Expect int64 for conversationID")
	}
	return id
}

func DeleteConversation(conversationId int64) {
	conversations, err := db.ListConversations()
	if err != nil {
		logger.PanicError(err, "Error listing conversations")
	}
	var title string
	for _, conv := range conversations {
		if conv.ID == conversationId {
			title = conv.Title
			break
		} else {
			logger.PanicError(err, "ConversationId provided does not match an ID in our database.\nYou can run the command without the --id flag to receive a prompt.")
		}
	}
	db.DeleteConversation(conversationId)
	fmt.Printf("Successfully deleted conversation.\nDeleted: %d - %s", conversationId, title)
}
