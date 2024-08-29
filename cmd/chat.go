package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/christianhturner/go-claude/chat"
	"github.com/christianhturner/go-claude/claude"
	cliui "github.com/christianhturner/go-claude/cli-ui"
	"github.com/christianhturner/go-claude/config"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	userMessage    string
	conversationId int64
	showHistory    bool
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Claude AI",
	Long: `This command allows you to chat with Claude AI. You can either provide a message and 
    conversationId directly using the --message and --id flag or enter a message when prompted.`,
	Run: func(cmd *cobra.Command, args []string) {
		model := viper.GetString(config.ModelKey)
		MaxTokens := viper.GetInt(config.MaxTokensKey)
		apiKey := viper.GetString(config.AnthropicApiKeyKey)
		stream := viper.GetBool(config.StreamKey)
		c := claude.NewClient(apiKey)

		convs, err := db.ListConversations()
		if err != nil {
			logger.PanicError(err, "Error listing conversations from DB")
		}
		if len(convs) == 0 {
			fmt.Println("No Conversations found.\nLet's get one created for you!\n")
			runCreateConversation()
		}

		if conversationId == 0 {
			conversationId = cliui.PromptForConversationId()
		}

		if showHistory && userMessage == "" {
			promptShowHistory, err := terminal.New().PromptConfirm("Would you like to see our conversation?")
			if err != nil {
				logger.PanicError(err, "Error getting user input.")
			}
			if promptShowHistory {
				displayAmount, messagePairs := cliui.PromptHistoricMessagePairs(conversationId)
				cliui.PresentHistoricMessagePairs(displayAmount, messagePairs)
			}
		}

		if userMessage == "" {
			userMessage = cliui.PromptUserForMessage()
		}

		history := chat.GetConversationHistory(conversationId)

		messageRequest := chat.MessageToRequest(userMessage)

		messages := chat.AppendHistoryToMessageRequest(messageRequest, history)

		ctx := context.Background()

		if stream {
			requestBody := claude.RequestBody{
				Model:     model,
				MaxTokens: MaxTokens,
				Messages:  messages,
				Stream:    stream,
			}
			stream := chat.StreamMessagesToClaude(ctx, requestBody, *c)
			var finalMessage claude.RequestMessages
			finalMessage.Role = claude.MessageRoleAssistant

			var contentBuilder strings.Builder
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			defer stream.Close()
			for {
				res, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					logger.PanicError(err, "Error with stream")
				}
				fmt.Printf("%s", res.Content[0].Text)
				contentBuilder.WriteString(res.Content[0].Text)
			}
			finalMessage.Content = contentBuilder.String()
			chat.AddMessageToConversationTable(conversationId, messageRequest)
			chat.AddMessageToConversationTable(conversationId, finalMessage)

		} else {

			requestBody := claude.RequestBody{
				Model:     model,
				MaxTokens: MaxTokens,
				Messages:  messages,
			}

			response := chat.SendMessageToClaude(ctx, requestBody, *c)

			fmt.Printf("Claude: %s/n", response.Content[0].Text)

			chat.AddMessageToConversationTable(conversationId, messageRequest)
			chat.AddMessageToConversationTable(conversationId, claude.RequestMessages{
				Role:    claude.MessageRoleAssistant,
				Content: response.Content[0].Text,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVar(&userMessage, "message", "", "Send a message to Claude")
	chatCmd.Flags().Int64Var(&conversationId, "id", 0, "Specify a Conversation by it's ID")
	chatCmd.Flags().BoolVarP(&showHistory, "history", "H", true, "Specify whether you want to see your last messages")
}
