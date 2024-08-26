package chat

import (
	"context"
	"fmt"
	"strconv"

	"github.com/christianhturner/go-claude/claude"
	"github.com/christianhturner/go-claude/config"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
	"github.com/spf13/viper"
)

var (
	convID    int64
	model     = viper.GetString(config.ModelKey)
	MaxTokens = viper.GetInt(config.MaxTokensKey)
	apiKey    = viper.GetString(config.AnthropicApiKeyKey)
)

func ChatWithoutFlags(message string) error {
	// c := claude.NewClient(apiKey)
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

	return nil
}

func ChatWithConvId(message, conversationId string) error {
	c := claude.NewClient(apiKey)
	var input string
	var err error

	if conversationId == "" {
		//	showConvList()
		fmt.Println("\n")
		conv, err := db.ListConversations()
		if err != nil {
			logger.PanicError(err, "Error Listing Conversations")
		}
		var messageIdOptions []string
		for _, convOptions := range conv {
			messageIdOptions = append(messageIdOptions, convOptions.Title)
		}
		id, input, err := terminal.New().PromptSelect("Select a conversation", messageIdOptions)
		if err != nil {
			fmt.Errorf("An Error occurred %v", err)
		}
		fmt.Printf("User selected %d: %s", id, input)
		convID = int64(id + 1)
	} else {
		idInt, err := strconv.Atoi(conversationId)
		if err != nil {
			logger.FatalError(err, "Please Provide a number for the messageID")
		}
		convID = int64(idInt + 1)
	}

	messages, err := db.GetMessages(convID)
	if err != nil {
		logger.FatalError(err, "Error getting messages with Id supplied")
	}
	logger.Info("\n Messages: \n", messages)
	var historicMessages []claude.RequestMessages
	for _, historicMessage := range messages {
		claudeMessage := claude.RequestMessages{
			Role:    historicMessage.Role,
			Content: historicMessage.Content,
		}
		historicMessages = append(historicMessages, claudeMessage)
	}

	if message == "" {
		input, err = terminal.New().Prompt("Start message to Claude: \n")
		if err != nil {
			logger.FatalError(err, "Error getting prompt")
		}
	} else {
		input = message
	}

	userMessage := claude.RequestMessages{
		Role:    claude.MessageRoleUser,
		Content: input,
	}
	finalMessages := append(historicMessages, userMessage)

	m := claude.RequestBody{
		Model:     model,
		MaxTokens: MaxTokens,
		Messages:  finalMessages,
	}

	ctx := context.Background()
	res, err := c.CreateMessages(ctx, m)
	if err != nil {
		logger.PanicError(err, "Panic on message")
	}
	fmt.Println(res.Content[0].Text)
	db.AddMessage(convID, claude.MessageRoleUser, userMessage.Content)
	db.AddMessage(convID, claude.MessageRoleAssistant, res.Content[0].Text)
	checkUpdate, err := db.GetConversation(convID)
	if err != nil {
		logger.PanicError(err, "Error with getting conversation")
	}
	logger.Info("\nUpdated Message:\n", checkUpdate)

	return nil
}
