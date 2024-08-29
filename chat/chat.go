package chat

import (
	"context"

	"github.com/christianhturner/go-claude/claude"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
)

type MessagePair struct {
	UserMessage      claude.RequestMessages
	AssistantMessage claude.RequestMessages
}

func GetMessagePairs(history []claude.RequestMessages) []MessagePair {
	pairs := []MessagePair{}
	for i := 0; i < len(history); i += 2 {
		if i+1 < len(history) {
			pair := MessagePair{
				UserMessage: claude.RequestMessages{
					Role:    claude.MessageRoleUser,
					Content: history[i].Content,
				},
				AssistantMessage: claude.RequestMessages{
					Role:    claude.MessageRoleAssistant,
					Content: history[i+1].Content,
				},
			}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

func GetConversationHistory(convId int64) []claude.RequestMessages {
	messages, err := db.GetMessages(convId)
	if err != nil {
		logger.PanicError(err, "Error getting messages from conversation table")
	}

	var historicMessages []claude.RequestMessages
	for _, historicMessage := range messages {
		claudeMessage := claude.RequestMessages{
			Role:    historicMessage.Role,
			Content: historicMessage.Content,
		}
		historicMessages = append(historicMessages, claudeMessage)
	}
	return historicMessages
}

func AppendHistoryToMessageRequest(messageRequest claude.RequestMessages, history []claude.RequestMessages) []claude.RequestMessages {
	return append(history, messageRequest)
}

func MessageToRequest(message string) claude.RequestMessages {
	return claude.RequestMessages{
		Role:    claude.MessageRoleUser,
		Content: message,
	}
}

func AddMessageToConversationTable(convId int64, message claude.RequestMessages) {
	db.AddMessage(convId, message.Role, message.Content)
}

func SendMessageToClaude(ctx context.Context, body claude.RequestBody, client claude.Client) *claude.ResponseBody {
	res, err := client.CreateMessages(ctx, body)
	if err != nil {
		logger.PanicError(err, "Panic sending message to claude")
	}
	return res
}

func StreamMessagesToClaude(ctx context.Context, body claude.RequestBody, client claude.Client) *claude.CreateMessagesStream {
	stream, err := client.CreateMessagesStream(ctx, body)
	if err != nil {
		logger.PanicError(err, "Error creating stream")
	}
	return stream
}
