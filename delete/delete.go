package delete

import (
	"fmt"

	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
)

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
