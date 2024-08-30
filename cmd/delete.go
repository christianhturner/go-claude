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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete subcommand will allow you to delete conversations and messages.",
	Long: `
    Use the delete subcommand to delete items from your stored data within go-claude.
    You can delete Conversations and Messages that are stored. These commands will allow you
    to delete items directly by specifying the subcommand followed by appropriate flags.
    The delete command cannot be used without a subcommand.

    go-claude delete conversation --[flags]

    From the conversation subcommand you can be provided a prompt if you supply no flags to specify which
    conversation that you wish to delete.

    go-claude delete conversation -> This will provide a prompt to select from your conversation list
    go-claude delete conversation --id 1 -> Will begin deleting the specified conversation with the Id provided, but will prompt if you're sure.
    go-claude delete conversation --id 1 -y -> This will directly delete conversation ID: 1 (excluding the -y will prompt if you're sure you wish to delete)

    go-claude delete messages --[flags]


    It's recommended to use the prompts to delete messages. The architecture of the database stores all messages in a single table, each message has a 
    conversation_id field, which is what sorts the messages with a conversation. The id used to delete a message via flags, is the id field, the messages
    actual Id. Therefore, you can delete messages across conversations with this command. The other issue, that isn't perfectly resolved with the prompt 
    in it's current state, is messages should be deleted in pairs. To continue a conversation, the historical messages must always be provided in pairs. 
    Every message from a user requires  a message from the assistant. Therefore, if you delete a message, without deleting it's pair, you'll break the 
    conversation. WIP to provide a better workflow.

    Example:
    go-claude delete messages -> This will provide a prompt to select a conversation, followed by another prompt to select messages you wish to delete.
    go-claude delete messages --id 1 -> This will take you to the prompt to select messages you wish to delete, but will bypass the requirement to select a conversation.
    go-claude delete messages --messId 8 -> This will directly delete the message with the id of 8. It will present a preview of the message prior to deleting.
    go-claude delete messages --messIds "1-3, 5, 7" -> This will delete messages with the ID 1, 2, 3, 5, and 7. It will present a preview of the messages prior to deleting.
    go-claude delete messages --messId 8 -y -> The '-y' flag will bypass the confirmation as long as the flags have been satisfied.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please provide a subcommand of conversation or message to define what you wish to delete.")
		cmd.Help()
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

// go-claude delete messages
var deleteMessages = &cobra.Command{
	Use:   "messages",
	Short: "Delete a messages from a specified conversation.",
	Long: `From the messages subcommand you can be provided a prompt if you supply no flags to specify which
    conversation that you wish to delete. 

    It's recommended to use the prompts to delete messages. The architecture of the database stores all messages in a single table, each message has a 
    conversation_id field, which is what sorts the messages with a conversation. The id used to delete a message via flags, is the id field, the messages
    actual Id. Therefore, you can delete messages across conversations with this command. The other issue, that isn't perfectly resolved with the prompt 
    in it's current state, is messages should be deleted in pairs. To continue a conversation, the historical messages must always be provided in pairs. 
    Every message from a user requires  a message from the assistant. Therefore, if you delete a message, without deleting it's pair, you'll break the 
    conversation. WIP to provide a better workflow.

    Example:
    go-claude delete messages -> This will provide a prompt to select a conversation, followed by another prompt to select messages you wish to delete.
    go-claude delete messages --id 1 -> This will take you to the prompt to select messages you wish to delete, but will bypass the requirement to select a conversation.
    go-claude delete messages --messId 8 -> This will directly delete the message with the id of 8. It will present a preview of the message prior to deleting.
    go-claude delete messages --messIds "1-3, 5, 7" -> This will delete messages with the ID 1, 2, 3, 5, and 7. It will present a preview of the messages prior to deleting.
    go-claude delete messages --messId 8 -y -> The '-y' flag will bypass the confirmation as long as the flags have been satisfied.`,
	Run: func(cmd *cobra.Command, args []string) {
		cliDeleteMessages(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteConversation, deleteMessages)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cliDeleteMessages(cmd *cobra.Command, args []string) {
	confirm, _ := cmd.Flags().GetBool("yes")
	messageIdIsSet := cmd.Flags().Changed("messId")
	messageIdsIsSet := cmd.Flags().Changed("messIds")
	// Does not support both values being set.
	if messageIdIsSet && messageIdsIsSet {
		fmt.Println("Please specify either --messId or --messIds, do not provide both")
		return
	}
	// If neither are set provide prompts.
	if !messageIdsIsSet && !messageIdIsSet {
		if conversationId == 0 {
			conversationId = cliui.PromptForConversationId()
		}
		messageIds := cliui.PromptMultiSelectMessageIds(conversationId)
		if confirm {
			for _, messId := range messageIds {
				delete.DeleteMessages(messId)
			}
		} else {
			userConfirmation := cliui.PromptForBool("Are you sure you want to delete:\n[%v]", messageIds)
			if userConfirmation {
				for _, messId := range messageIds {
					delete.DeleteMessages(messId)
				}
			}
		}

	}
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
				err := delete.DeleteConversation(conversationId)
				if err != nil {
					logger.PanicError(err, "ConversationId provided does not match an ID in our database.\nYou can run the command without the --id flag to receive a prompt.")
				}
				fmt.Printf("Successfully deleted conversation.\nDeleted: %d - %s", conversationId, title)
			} else {
				fmt.Println("Cancelling request to delete Conversation.")
			}
		}
	}
}
