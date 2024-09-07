/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	cliui "github.com/christianhturner/go-claude/cli-ui"
	"github.com/christianhturner/go-claude/go_claude_list"
	"github.com/christianhturner/go-claude/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listConversationsCmd, listMessagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list allows you to list go-claude data.",
	Long: `list various items by following this command with a supported subcommand. You can
    list conversations, messages, global configuration, and conversation level configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide a subcommand [conversations, messages, global-configuration, or conversation-configuration]")
	},
}

// go-claude list conversations
var listConversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "List conversations in a table view",
	Long:  `List conversations into a table view that provides columns including ID, Title, Created, and Last Updated`,
	Run: func(cmd *cobra.Command, args []string) {
		cliRunListConversationsFunction(cmd, args)
	},
}

var listMessagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "List messages in a table view",
	Long:  "List messages in a table view that provides columns which include Id, Role, Content, and Created",
	Run: func(cmd *cobra.Command, args []string) {
		cliRunListMessagesFunction(cmd, args)
	},
}

// go-claude

func cliRunListConversationsFunction(cmd *cobra.Command, args []string) {
	conv, err := go_claude_list.GetAllConversations()
	if err != nil {
		logger.FatalError(err, "Error getting conversations")
	}
	cliui.ShowConvList(conv)
}

func cliRunListMessagesFunction(cmd *cobra.Command, args []string) {
	conversationId := cliui.PromptForConversationId()
	go_claude_list.ShowMessageList(conversationId)
}
