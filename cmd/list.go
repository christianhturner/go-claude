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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		showConvList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func showConvList() error {
	conv, err := db.ListConversations()
	if err != nil {
		logger.FatalError(err, "Error listing conversations.")
	}

	term := terminal.New()

	table := term.NewTable(20)

	idMaxWidth := 7
	table.AddColumn("ID", "ID", 5, &idMaxWidth, false, 0)
	table.AddColumn("Title", "Title", 40, nil, true, 0)
	timeMaxWidth := 19
	table.AddColumn("Created", "Created", 19, &timeMaxWidth, false, terminal.AlignCenter)
	table.AddColumn("Updated", "Updated", 19, &timeMaxWidth, false, terminal.AlignCenter)

	for _, c := range conv {
		table.AddRow(map[string]interface{}{
			"ID":        c.ID,
			"Title":     c.Title,
			"Something": c.Title,
			"Created":   c.CreatedAt.Format("2006-01-02 15:04:05"),
			"Updated":   c.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	table.Render()

	return nil
}
