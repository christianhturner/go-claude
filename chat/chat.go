package chat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/christianhturner/go-claude/claude"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
)

type MessagePair struct {
	UserMessage      claude.RequestMessages
	AssistantMessage claude.RequestMessages
}

func getMessagePairs(history []claude.RequestMessages) []MessagePair {
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

func PrintMessageHistory(conversationId int64) {
	history := GetConversationHistory(conversationId)
	pairs := getMessagePairs(history)

	printHistory, err := terminal.New().PromptConfirm("Would you like to see our conversation?")
	if err != nil {
		logger.PanicError(err, "Error prompting for message history.")
	}

	if !printHistory {
		return
	}

	maxPairs := len(pairs)
	var numPairs int

	for {
		input, err := terminal.New().Prompt(fmt.Sprintf("How many message pairs would you like to review? (1-%d)", maxPairs))
		if err != nil {
			logger.PanicError(err, "Error reading user input.")
			return
		}
		numPairs, err = strconv.Atoi(input)
		if err != nil || numPairs < 1 || numPairs > maxPairs {
			fmt.Printf("Please enter a valid number between 1 and %d.\n", maxPairs)
			continue
		}
		break
	}

	fmt.Println("\nConversation History:")
	// Change this loop to iterate from the beginning
	for i := 0; i < numPairs && i < len(pairs); i++ {
		pair := pairs[i]
		fmt.Printf("\nUser: %s\n", pair.UserMessage.Content)
		fmt.Printf("\nClaude: %s\n", pair.AssistantMessage.Content)
	}
	// Extra line space
	fmt.Println("\n")
}

func PromptUserForMessage() string {
	input, err := terminal.New().Prompt("User: ")
	if err != nil {
		logger.PanicError(err, "Error prompting user for message.")
	}
	return input
}

func PromptUserForConversationId() int64 {
	conv, err := db.ListConversations()
	if err != nil {
		logger.PanicError(err, "Error listing conversations")
	}
	options := make(map[interface{}]string)
	for _, convOptions := range conv {
		options[convOptions.ID] = convOptions.Title
	}
	selected := terminal.New().PromptOptionsSelect(options)
	fmt.Printf("Selected: ID=%v, Description=%s\n", selected.ID, selected.Description)

	id, ok := selected.ID.(int64)
	if !ok {
		logger.PanicError(err, "Invalid ID type; Expected int64 for conversation ID")
	}

	return id
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
