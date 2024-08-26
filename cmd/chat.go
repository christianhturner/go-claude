package cmd

import (
	"github.com/christianhturner/go-claude/chat"
	"github.com/spf13/cobra"
)

var (
	userMessage    string
	conversationId string
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Claude AI",
	Long: `This command allows you to chat with Claude AI. You can either provide a message
directly using the --message flag or enter a message when prompted.`,
	Run: func(cmd *cobra.Command, args []string) {
		// chat.ChatWithConvId(userMessage, conversationId)
		chat.ChatWithoutFlags("")
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringVar(&userMessage, "message", "", "Send a message to Claude")
	chatCmd.Flags().StringVar(&conversationId, "id", "", "Specify a messageId")
}
