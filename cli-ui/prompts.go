package cliui

import (
	"fmt"
	"strconv"

	"github.com/christianhturner/go-claude/chat"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
)

func PromptHistoricMessagePairs(conversationId int64) (int, []chat.MessagePair) {
	history := chat.GetConversationHistory(conversationId)
	messagePairs := chat.GetMessagePairs(history)

	maxPairs := len(messagePairs)
	if maxPairs == 0 {
		fmt.Println("No historic messages to print.")
		return 0, nil
	}
	var numPairs int

	for {
		input, err := terminal.New().Prompt(fmt.Sprintf("How many message pairs would you like to review? (1-%d)", maxPairs))
		if err != nil {
			logger.PanicError(err, "Error reading user input.")
		}
		numPairs, err = strconv.Atoi(input)
		if err != nil || numPairs < 1 || numPairs > maxPairs {
			fmt.Printf("Please enter a valid number between 1 and %d.\n", maxPairs)
			continue
		}
		break
	}
	return numPairs, messagePairs
}

func PromptUserForMessage() string {
	input, err := terminal.New().Prompt("User: ")
	if err != nil {
		logger.PanicError(err, "Error prompting user for message.")
	}
	return input
}

func PromptForBool(format string, a ...any) bool {
	userConfirmation, err := terminal.New().PromptConfirm(fmt.Sprintf(format, a...))
	if err != nil {
		logger.PanicError(err, "Error getting user confirmation")
	}
	return userConfirmation
}

func PromptMultiSelectMessageIds(conversationId int64) []int64 {
	messages, err := db.GetMessages(conversationId)
	if err != nil {
		logger.PanicError(err, "Error listing conversations")
	}
	options := make(map[interface{}]string)
	for _, messOptions := range messages {
		options[messOptions.ID] = messOptions.Content
	}
	selectedOptions := terminal.New().PromptMultipleOptionsSelect(options)
	var messageIds []int64
	fmt.Println("\nSelected:\n")
	for _, messIds := range selectedOptions {
		id, ok := messIds.ID.(int64)
		if !ok {
			logger.PanicError(err, "ID is not int64")
		}
		fmt.Printf("\n%d - %s\n", id, messIds.Description)
		messageIds = append(messageIds, id)
	}
	return messageIds
}

func PromptForConversationId() int64 {
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
