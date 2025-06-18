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
	var buf bytes.Buffer
	originalLogger := logger
	logger = log.New(&buf, "", 0)
	defer func() { logger = originalLogger }()

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

	SetLogLevel("debug")
	buf.Reset()
	Debug("debug message")
	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected log to contain 'debug message', got '%s'", output)
	}
}

func TestEnvironmentVariable(t *testing.T) {
	originalEnv := os.Getenv("LOG_LEVEL")
	defer os.Setenv("LOG_LEVEL", originalEnv)

	os.Setenv("LOG_LEVEL", "error")

	currentLevel = INFO

	var buf bytes.Buffer
	originalLogger := logger
	logger = log.New(&buf, "", 0)
	defer func() { logger = originalLogger }()

	Init()

	if currentLevel != ERROR {
		t.Errorf("Expected log level ERROR, got %v", currentLevel)
	}

	buf.Reset()
	Info("info message")
	if buf.String() != "" {
		t.Errorf("Expected no log output for INFO message, got '%s'", buf.String())
	}

	buf.Reset()
	Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Errorf("Expected log to contain 'error message', got '%s'", buf.String())
	}
}
