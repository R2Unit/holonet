package logger

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"info", INFO},
		{"warn", WARN},
		{"error", ERROR},
		{"fatal", FATAL},
		{"invalid", INFO},
	}

	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			SetLogLevel(test.level)
			if currentLevel != test.expected {
				t.Errorf("Expected log level %v, got %v", test.expected, currentLevel)
			}
		})
	}
}

func TestLogOutput(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	originalLogger := logger
	logger = log.New(&buf, "", 0) // Remove timestamp for easier testing
	defer func() { logger = originalLogger }()

	// Test each log level
	tests := []struct {
		level       LogLevel
		logFunc     func(string, ...interface{})
		message     string
		shouldLog   bool
		levelPrefix string
	}{
		{INFO, Info, "info message", true, "[INFO]"},
		{ERROR, Error, "error message", true, "[ERROR]"},
		{DEBUG, Debug, "debug message", false, "[DEBUG]"}, // Shouldn't log at INFO level
		{WARN, Warn, "warn message", true, "[WARN]"},
	}

	// Set log level to INFO
	SetLogLevel("info")

	for _, test := range tests {
		buf.Reset()
		test.logFunc(test.message)
		output := buf.String()

		if test.shouldLog {
			if !strings.Contains(output, test.message) {
				t.Errorf("Expected log to contain message '%s', got '%s'", test.message, output)
			}
			if !strings.Contains(output, test.levelPrefix) {
				t.Errorf("Expected log to contain prefix '%s', got '%s'", test.levelPrefix, output)
			}
		} else {
			if output != "" {
				t.Errorf("Expected no log output, got '%s'", output)
			}
		}
	}

	// Set log level to DEBUG and test again
	SetLogLevel("debug")
	buf.Reset()
	Debug("debug message")
	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected log to contain 'debug message', got '%s'", output)
	}
}

func TestEnvironmentVariable(t *testing.T) {
	// Save original environment and restore it after the test
	originalEnv := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", originalEnv)

	// Set environment variable
	os.Setenv("LOG_LEVEL", "error")

	// Reset currentLevel
	currentLevel = INFO

	// Capture log output
	var buf bytes.Buffer
	originalLogger := logger
	logger = log.New(&buf, "", 0)
	defer func() { logger = originalLogger }()

	// Initialize logger (should read from environment)
	Init()

	// Check that log level was set correctly
	if currentLevel != ERROR {
		t.Errorf("Expected log level ERROR, got %v", currentLevel)
	}

	// Test that INFO messages are not logged
	buf.Reset()
	Info("info message")
	if buf.String() != "" {
		t.Errorf("Expected no log output for INFO message, got '%s'", buf.String())
	}

	// Test that ERROR messages are logged
	buf.Reset()
	Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Errorf("Expected log to contain 'error message', got '%s'", buf.String())
	}
}
