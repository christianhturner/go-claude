/*
Copyright © 2024 Christian Turner <ch.turner94@gmail.com | github.com/christianhturner>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/christianhturner/go-claude/config"
	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
	"github.com/spf13/cobra"
)

var stopChan = make(chan os.Signal, 1)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-claude",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig, logger.InitLogger, initDB, sessionInit)
	config.AddFlags(rootCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sessionInit() {
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		fmt.Println("\nShutting Down...")
		logger.Debug("\nShutting Down...")
		db.Close()
		os.Exit(0)
	}()
}

func initDB() {
	home, err := os.UserHomeDir()
	logger.FatalError(err, "Failed to get user home directory")

	configDir := filepath.Join(home, ".config", "go-claude")
	err = os.MkdirAll(configDir, 0755)
	logger.FatalError(err, "Failed to create config directory")

	dbPath := filepath.Join(configDir, "data.db")
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		logger.LogError(err, "Failed to create database file")
	}
	err = db.InitDatabase(dbPath)
	logger.FatalError(err, "Failed to initialize database")
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := os.UserHomeDir()
// 		cobra.CheckErr(err)
//
// 		configDir := filepath.Join(home, ".config", "go-claude")
// 		err = os.MkdirAll(configDir, 0755)
// 		cobra.CheckErr(err)
//
// 		configName := "config"
// 		configType := "json"
// 		configPath := filepath.Join(configDir, configName+"."+configType)
//
// 		// Search config in home directory with name ".go-claude" (without extension).
// 		viper.AddConfigPath(configDir)
// 		viper.SetConfigName(configName)
// 		viper.SetConfigType(configType)
//
// 		_, err = os.Stat(configPath)
// 		if !os.IsExist(err) {
// 			SetDefaults()
// 			viper.SafeWriteConfig()
// 		}
//
// 	}
//
// 	viper.AutomaticEnv() // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	err := viper.ReadInConfig()
// 	if err == nil {
// 		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
// 	} else {
// 		fmt.Errorf("Failed to read config file: %w", err)
// 	}
//
// 	viper.OnConfigChange(func(e fsnotify.Event) {
// 		fmt.Printf("Config file changed", e.Name)
// 	})
// 	viper.WatchConfig()
//
// 	checkApiKey := viper.GetString("Anthropic_API_Key")
// 	if checkApiKey == "" {
// 		term := terminal.New()
// 		userInput, err := term.Prompt("Please provide your Anthroipic API Key:\n")
// 		if err != nil {
// 			fmt.Errorf("Error requesting user input for API key: %v", err)
// 		}
// 		viper.Set("Anthropic_API_Key", userInput)
// 		viper.WriteConfig()
// 	}
// }
