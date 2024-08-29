/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	cliui "github.com/christianhturner/go-claude/cli-ui"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/delete"
	"github.com/christianhturner/go-claude/logger"
	"github.com/spf13/cobra"
)

var confirm bool

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete subcommand will allow you to delete conversations and messages.",
	Long: `Use the delete subcommand to delete items from your stored data within go-claude.
    You can delete Conversations and Messages that are stored. These commands will allow you
    to delete items directly by specifying the subcommand followed by appropriate flags.
    The delete command cannot be used without a following subcommand. Example:

    go-claude delete conversation -> This will provide a prompt to select from your conversation list
    go-claude delete message -> This will provide a prompt to select first a conversation, followed by a prompt to choose a message.

    go-claude delete conversation --id 1 -y -> This will directly delete conversation ID: 1 (excluding the -y will prompt if you're sure you wish to delete)
    go-claude delete message --id 1 --messageId 23 -y -> This will directly delete Message ID: 23 from conversation ID: 1 (excluding the -y will prompt if you're sure you wish to delete)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide a subcommand of conversation or message to define what you wish to delete.")
	},
}

// go-claude delete conversation
var deleteConversation = &cobra.Command{
	Use:   "conversations",
	Short: "Delete a converastion from your Claude conversation list.",
	Long: `From the conversation subcommand you can be provided a prompt if you supply no flags to specify which
    conversation that you wish to delete. Example:
    go-claude delete conversation -> This will provide a prompt to select from your conversation list
    go-claude delete conversation --id 1 -> Will begin deleting the specified conversation with the Id provided, but will prompt if you're sure.
    go-claude delete conversation --id 1 -y -> This will directly delete conversation ID: 1 (excluding the -y will prompt if you're sure you wish to delete)`,
	Run: func(cmd *cobra.Command, args []string) {
		cliDelteConversation(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteConversation)
	deleteConversation.Flags().Int64Var(&conversationId, "id", 0, "Specify a Conversation by it's ID")
	deleteConversation.Flags().BoolP("yes", "y", false, "Automatically confirm without prompts.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cliDelteConversation(cmd *cobra.Command, args []string) {
	confirm, _ := cmd.Flags().GetBool("yes")
	if conversationId == 0 {
		conversationId = cliui.PromptForConversationId()
		delete.DeleteConversation(conversationId)
	} else {
		if confirm {
			delete.DeleteConversation(conversationId)
		} else {
			conversations, err := db.ListConversations()
			var title string
			for _, conv := range conversations {
				if conv.ID == conversationId {
					title = conv.Title
					break
				} else {
					logger.PanicError(err, "ConversationId provided does not match an Id in our database.\nYou can run the command without the --id flag to receive a prompt.")
				}
			}
			userConfirmation := cliui.PromptForBool("Are you sure you want to delete:\nID: %v - %v", conversationId, title)
			if userConfirmation {
				delete.DeleteConversation(conversationId)
			} else {
				fmt.Println("Cancelling request to delete Conversation.")
			}
		}
	}
}
