/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/christianhturner/go-claude/config"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Global configurations for go-claude.",
	Long: `Cofigure command will direclty configure your gloabl
    configuration. Any configuration can be set through this command
    followed by a '--flag' followed by the value you wish to set.

    This will change the value within the global config file. The
    default location of this file is ~/.config/go-claude/config.json, 
    though this location can be changed within the config file and 
    configure command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("defaults").Changed {
			config.ResetToDefaults()
			fmt.Println("configuration reset to defaults")
		} else {
			config.UpdateConfig(cmd)
			fmt.Println("Configuration file updated")
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().Bool("defaults", false, "Reset configuration to default values.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
