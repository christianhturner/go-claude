/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/christianhturner/go-claude/terminal"
	"github.com/spf13/cobra"
)

var conversationTitle string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			err := runCreate()
			if err != nil {
				logger.FatalError(err, "Error executing runCreate")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&conversationTitle, "title", "", "Title for conversation.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCreate() error {
	term := terminal.New()
	userSelect, err := term.PromptConfirm("Would you like to give your conversation a name?")
	logger.LogError(err, "Error making a selection at runCreate.")
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
