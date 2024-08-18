package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

// InitLogger initializes the global logger instance.
// It should be called once at the beginning of the program.
func InitLogger() {
	logger, err := initZap()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	sugar = logger.Sugar()
}

// Debug logs the provided arguments at \[DebugLevel].
// Spaces are added between arguments when neither is a string.
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Info logs the provided arguments at \[InfoLevel].
// Spaces are added between arguments when neither is a string.
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// Warn logs the provided arguments at \[WarnLevel].
// Spaces are added between arguments when neither is a string.
func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

// Error logs the provided arguments at \[ErrorLevel].
// Spaces are added between arguments when neither is a string.
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// Fatal constructs a message with the provided arguments and calls os.Exit.
// Spaces are added between arguments when neither is a string.
func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
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
