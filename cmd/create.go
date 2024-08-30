/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	cliui "github.com/christianhturner/go-claude/cli-ui"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Use this command to create new conversations.",
	Long: `Use this command to create new conversations. Most likely subcommands will be
    added, and create will not work without the use of the additional subcommands. Currently,
    no other items to be created, so leaving this as is.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("title") {
			if conversationTitle == "none" {
				id, err := db.CreateConversation("")
				if err != nil {
					logger.FatalError(err, "Error creating conversation with no title at go-claude create --title \"title\"")
				}
				logger.Debug("Created a conversation with no name. Id: ", id)
			} else {
				id, err := db.CreateConversation(conversationTitle)
				if err != nil {
					logger.FatalError(err, "Error creating conversation with title")
				}
				logger.Debug("Created conversation: \nId: ", id, ": Title", conversationTitle)
			}
		} else {
			err := runCreateConversation()
			if err != nil {
				logger.FatalError(err, "Error executing runCreate")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCreateConversation() error {
	term := terminal.New()
	userSelect := cliui.PromptForBool("Would you like to give your conversation a name?", nil)
	switch userSelect {
	case true:
		input, err := term.Prompt("Please provide a name for your conversation:")
		logger.LogError(err, "Error inputting name at runCreate")
		id, err := db.CreateConversation(input)
		logger.FatalError(err, "Error creating conversation in database at runCreate")
		logger.Debug("Created a conversation:\nId:", id, ": Title: ", input)
	case false:
		id, err := db.CreateConversation("")
		logger.FatalError(err, "Error creating conversation with no title at runCreate")
		logger.Debug("Created a conversation with no name. Id: ", id)
	}
	return nil
}
