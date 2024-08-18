/*
Copyright © 2024 Christian Turner <ch.turner94@gmail.com | github.com/christianhturner>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/christianhturner/go-claude/pkg/db"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfgFile string
	log     *zap.SugaredLogger
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
	cobra.OnInitialize(initConfig, initLogger, initDB)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.go-claude/config.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initDB() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configDir := filepath.Join(home, ".config", "go-claude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		cobra.CheckErr(err)
	}
	dbPath := filepath.Join(configDir, "data.db")
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		_, err := os.Create(dbPath)
		if err != nil {
			fmt.Printf("failed to create file %v", err)
		}
	}
	db.InitDatabase(dbPath)
}

func initLogger() {
	logger, err := initZap()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log = logger.Sugar()
	defer log.Sync()
}

func initZap() (*zap.Logger, error) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configDir := filepath.Join(home, ".config", "go-claude")
	logFile := filepath.Join(configDir, "go-claude.log")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Errorf("failed to create config directory: %w", err)
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Errorf("failed to open log file: %w", err)
	}

	logLevel := viper.GetString("log_level")
	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	logger := zap.New(core)
	return logger, nil
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
		if err := os.MkdirAll(configDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
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
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}
