package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/christianhturner/go-claude/pkg/claude"
	"github.com/christianhturner/go-claude/pkg/db"
	"github.com/christianhturner/go-claude/pkg/log"
	"github.com/christianhturner/go-claude/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	claudeMessage string
	messageId     string
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Claude AI",
	Long: `This command allows you to chat with Claude AI. You can either provide a message
directly using the --chat flag or enter a message when prompted.`,
	Run: func(cmd *cobra.Command, args []string) {
		runChat(claudeMessage, messageId)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVar(&claudeMessage, "chat", "", "Send a message to Claude")
	chatCmd.Flags().StringVar(&messageId, "id", "", "Specify a messageId")
}

func runChat(message, conversationId string) error {
	var convID int64
	const (
		model     = "claude-3-5-sonnet-20240620"
		MaxTokens = 2000
	)
	apiKey := viper.GetString("anthropic_api_key")

	c := claude.NewClient(apiKey)

	var input string
	var err error

	if conversationId == "" {
		//	showConvList()
		fmt.Println("\n")
		conv, err := db.ListConversations()
		if err != nil {
			log.PanicError(err, "Error Listing Conversations")
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
			log.FatalError(err, "Please Provide a number for the messageID")
		}
		convID = int64(idInt + 1)
	}

	messages, err := db.GetMessages(convID)
	if err != nil {
		log.FatalError(err, "Error getting messages with Id supplied")
	}
	log.Info("\n Messages: \n", messages)
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
			log.FatalError(err, "Error getting prompt")
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
		log.PanicError(err, "Panic on message")
	}
	fmt.Println(res.Content[0].Text)
	db.AddMessage(convID, claude.MessageRoleUser, userMessage.Content)
	db.AddMessage(convID, claude.MessageRoleAssistant, res.Content[0].Text)
	checkUpdate, err := db.GetConversation(convID)
	if err != nil {
		log.PanicError(err, "Error with getting conversation")
	}
	log.Info("\nUpdated Message:\n", checkUpdate)

	return nil
}
