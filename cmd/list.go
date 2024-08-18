/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/christianhturner/go-claude/pkg/db"
	"github.com/christianhturner/go-claude/pkg/log"
	"github.com/christianhturner/go-claude/pkg/terminal"
	"github.com/mattn/go-runewidth"
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
		runList()
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

func runList() error {
	conv, err := db.ListConversations()
	if err != nil {
		log.FatalError(err, "Error listing conversations.")
	}

	width, _ := terminal.GetWidthAndHeight()

	// Define fixed widths
	idWidth := 5
	createdWidth := 19
	updatedWidth := 19
	minTitleWidth := 20
	separatorWidth := 3 // For " | "

	// Calculate title width
	titleWidth := width - idWidth - updatedWidth - separatorWidth*3
	showCreated := true

	if titleWidth < minTitleWidth+createdWidth {
		// Drop 'Created' column if not enough space
		titleWidth += createdWidth + separatorWidth
		showCreated = false
	}

	// Print headers
	idHeader := centerString("ID", idWidth)
	titleHeader := centerString("Title", titleWidth)
	createdHeader := centerString("Created", createdWidth)
	updatedHeader := centerString("Updated", updatedWidth)

	if showCreated {
		fmt.Printf("\n%s | %s | %s | %s\n", idHeader, titleHeader, createdHeader, updatedHeader)
	} else {
		fmt.Printf("\n%s | %s | %s\n", idHeader, titleHeader, updatedHeader)
	}

	// Print separator line
	fmt.Println(strings.Repeat("-", width))

	// Print conversations
	for _, c := range conv {
		id := formatID(int(c.ID), idWidth)
		title := formatTitle(c.Title, titleWidth)
		created := c.CreatedAt.Format("2006-01-02 15:04:05")
		updated := c.UpdatedAt.Format("2006-01-02 15:04:05")

		if showCreated {
			fmt.Printf("%s | %s | %s | %s\n", id, title, created, updated)
		} else {
			fmt.Printf("%s | %s | %s\n", id, title, updated)
		}
	}

	return nil
}

func centerString(s string, width int) string {
	sWidth := runewidth.StringWidth(s)
	if sWidth >= width {
		return runewidth.Truncate(s, width, "")
	}
	leftPad := (width - sWidth) / 2
	rightPad := width - sWidth - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

func formatID(id int, width int) string {
	return fmt.Sprintf("%*d", width, id)
}

func formatTitle(title string, width int) string {
	if title == "" {
		title = "[No Title]"
	}
	if runewidth.StringWidth(title) <= width {
		return runewidth.FillRight(title, width)
	}
	truncated := runewidth.Truncate(title, width-3, "...")
	return runewidth.FillRight(truncated, width)
}
