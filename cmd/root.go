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
