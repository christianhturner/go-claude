package logger

import (
	"fmt"
	"os"
	"runtime/debug"
)

// FatalError logs a fatal error message if the error is not nil.
// It then calls os.Exit(1).
func FatalError(err error, message string) {
	if err != nil {
		sugar.Fatal("%s: %v", message, err)
	}
}

// PanicError logs an error message and panics if the error is not nil.
func PanicError(err error, message string) {
	if err != nil {
		sugar.Panicf("%s: %v", message, err)
	}
}

// LogError logs an error message if the error is not nil.
func LogError(err error, message string) {
	if err != nil {
		sugar.Errorf("%s: %v", message, err)
	}
}

// WarnError logs a warning message if the error is not nil.
func WarnError(err error, message string) {
	if err != nil {
		sugar.Warnf("%s: %v", message, err)
	}
}

// RecoverPanic recovers from panics, logs the error, and optionally exits the program.
// This function should be deferred at the beginning of main functions or goroutines.
func RecoverPanic(exit bool) {
	if r := recover(); r != nil {
		sugar.Errorf("Recovered from panic: %v\nStack trace:\n%s", r, debug.Stack())
		if exit {
			os.Exit(1)
		}
	}
}

// AssertNoError panics if the error is not nil.
// This is useful for errors that should never happen during normal operation.
func AssertNoError(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", message, err))
	}
}
