package logger

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// DEBUG level for detailed troubleshooting information
	DEBUG LogLevel = iota
	// INFO level for general operational information
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
	// FATAL level for critical errors that cause the program to exit
	FATAL
)

var (
	// currentLevel is the current log level
	currentLevel LogLevel = INFO
	logger                = log.New(os.Stdout, "", log.LstdFlags)
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

var levelPrefixes = map[LogLevel]string{
	DEBUG: "[DEBUG] ",
	INFO:  "[INFO] ",
	WARN:  "[WARN] ",
	ERROR: "[ERROR] ",
	FATAL: "[FATAL] ",
}

func Init() {
	logLevelFlag := flag.String("log-level", "info", "Set the logging level (debug, info, warn, error, fatal)")
	flag.Parse()
	logLevelEnv := os.Getenv("LOG_LEVEL")
	logLevel := *logLevelFlag
	if logLevelEnv != "" {
		logLevel = logLevelEnv
		Debug("Using log level from environment variable: %s", logLevelEnv)
	}

	SetLogLevel(logLevel)

	Info("Logger initialized with level: %s", levelNames[currentLevel])
}

func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	case "fatal":
		currentLevel = FATAL
	default:
		currentLevel = INFO
		logger.Printf("Unknown log level '%s', defaulting to INFO", level)
	}
}

func Debug(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		logger.Print(levelPrefixes[DEBUG] + fmt.Sprintf(format, v...))
	}
}

func Info(format string, v ...interface{}) {
	if currentLevel <= INFO {
		logger.Print(levelPrefixes[INFO] + fmt.Sprintf(format, v...))
	}
}

func Warn(format string, v ...interface{}) {
	if currentLevel <= WARN {
		logger.Print(levelPrefixes[WARN] + fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		logger.Print(levelPrefixes[ERROR] + fmt.Sprintf(format, v...))
	}
}

func Fatal(format string, v ...interface{}) {
	if currentLevel <= FATAL {
		logger.Fatal(levelPrefixes[FATAL] + fmt.Sprintf(format, v...))
	}
}

func GetLogLevel() LogLevel {
	return currentLevel
}

func GetLogLevelName() string {
	return levelNames[currentLevel]
}
