// Package logger provides a logging interface and implementations using the
// charmbracelet/lipgloss and charmbracelet/log packages. It includes a
// configurable logger (CharmLogger) and a test logger (NullLogger) for capturing
// log output in tests.
//
// Usage:
//
// Basic usage of CharmLogger:
//
//	package main
//
//	import (
//		"os"
//		"github.com/yourusername/go-common/pkg/tui/logger/charm"
//	)
//
//	func main() {
//		log := logger.New(os.Stdout, logger.LogLevelInfo)
//		log.Info("This is an info message")
//		log.Debug("This is a debug message") // This won't be printed because the level is set to Info
//	}
//
// Using NullLogger for testing:
//
//	package main
//
//	import (
//		"testing"
//		"github.com/yourusername/go-common/pkg/tui/logger/charm"
//	)
//
//	func TestLogging(t *testing.T) {
//		log := logger.NewTestLogger(t)
//		log.Info("This is an info message")
//		log.Debug("This is a debug message")
//		// The log output can be accessed via log.LogOutput.String()
//		if !strings.Contains(log.LogOutput.String(), "This is an info message") {
//			t.Error("Expected info message to be logged")
//		}
//	}
package logger
