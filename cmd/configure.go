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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

// func setAPIKey(apiKey string) error {
// }

// func SetDefaults() {
// 	checkApiKey := viper.IsSet("Anthropic_API_Key")
// 	if checkApiKey == false {
// 		viper.SetDefault("Anthropic_API_Key", "")
// 	}
// 	viper.SetDefault("log_level", "info")
// }
