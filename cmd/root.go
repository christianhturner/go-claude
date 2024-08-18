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

	"github.com/christianhturner/go-claude/pkg/db"
	"github.com/christianhturner/go-claude/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	stopChan = make(chan os.Signal, 1)
)

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
	cobra.OnInitialize(initConfig, log.InitLogger, initDB, sessionInit)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.go-claude/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sessionInit() {
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		log.Info("\nShutting Down...")
		db.Close()
		os.Exit(0)
	}()
}

func initDB() {
	home, err := os.UserHomeDir()
	log.FatalError(err, "Failed to get user home directory")

	configDir := filepath.Join(home, ".config", "go-claude")
	err = os.MkdirAll(configDir, 0755)
	log.FatalError(err, "Failed to create config directory")

	dbPath := filepath.Join(configDir, "data.db")
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		log.LogError(err, "Failed to create database file")
	}
	err = db.InitDatabase(dbPath)
	log.FatalError(err, "Failed to initialize database")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configDir := filepath.Join(home, ".config", "go-claude")
		err = os.MkdirAll(configDir, 0755)
		cobra.CheckErr(err)

		configName := "config"
		configType := "json"
		// configPath := filepath.Join(configDir, configName+"."+configType)

		// Search config in home directory with name ".go-claude" (without extension).
		viper.AddConfigPath(configDir)
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		viper.SafeWriteConfig()

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Printf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		fmt.Errorf("Failed to read config file: %w", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed", e.Name)
	})
	viper.WatchConfig()
}
