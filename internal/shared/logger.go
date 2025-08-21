package shared

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "ERROR"
	case LogLevelWarn:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging for the pipeline analyzer
type Logger struct {
	level      LogLevel
	fileLogger *log.Logger
	file       *os.File
}

var globalLogger *Logger

// InitLogger initializes the global logger
func InitLogger(level LogLevel, logDir string) error {
	if globalLogger != nil {
		globalLogger.Close()
	}

	logger := &Logger{
		level: level,
	}

	// Create log directory if it doesn't exist
	if logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// Create log file with timestamp
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logFile := filepath.Join(logDir, fmt.Sprintf("pipeline-analyzer_%s.log", timestamp))
		
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to create log file: %w", err)
		}

		logger.file = file
		logger.fileLogger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)
		
		// Log initialization
		logger.logToFile(LogLevelInfo, "Logger", "Logger initialized", map[string]interface{}{
			"log_file": logFile,
			"level":    level.String(),
		})
	}

	globalLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		// Initialize with default settings if not already initialized
		InitLogger(LogLevelInfo, "")
	}
	return globalLogger
}

// Close closes the logger and any open file handles
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
		l.file = nil
		l.fileLogger = nil
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// logToFile writes to the log file if file logging is enabled
func (l *Logger) logToFile(level LogLevel, component, message string, context map[string]interface{}) {
	if l.fileLogger == nil {
		return
	}

	// Format context as key=value pairs
	var contextStr strings.Builder
	for key, value := range context {
		if contextStr.Len() > 0 {
			contextStr.WriteString(" ")
		}
		contextStr.WriteString(fmt.Sprintf("%s=%v", key, value))
	}

	logEntry := fmt.Sprintf("[%s] [%s] %s", level.String(), component, message)
	if contextStr.Len() > 0 {
		logEntry += " " + contextStr.String()
	}

	l.fileLogger.Println(logEntry)
}

// shouldLog checks if a message should be logged at the given level
func (l *Logger) shouldLog(level LogLevel) bool {
	return level <= l.level
}

// Error logs an error message
func (l *Logger) Error(component, message string, context map[string]interface{}) {
	if l.shouldLog(LogLevelError) {
		if context == nil {
			context = make(map[string]interface{})
		}
		l.logToFile(LogLevelError, component, message, context)
		
		// Also print to stderr for errors
		fmt.Fprintf(os.Stderr, "ERROR [%s]: %s\n", component, message)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(component, message string, context map[string]interface{}) {
	if l.shouldLog(LogLevelWarn) {
		if context == nil {
			context = make(map[string]interface{})
		}
		l.logToFile(LogLevelWarn, component, message, context)
	}
}

// Info logs an info message
func (l *Logger) Info(component, message string, context map[string]interface{}) {
	if l.shouldLog(LogLevelInfo) {
		if context == nil {
			context = make(map[string]interface{})
		}
		l.logToFile(LogLevelInfo, component, message, context)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(component, message string, context map[string]interface{}) {
	if l.shouldLog(LogLevelDebug) {
		if context == nil {
			context = make(map[string]interface{})
		}
		l.logToFile(LogLevelDebug, component, message, context)
	}
}

// ParseError logs a parsing error with detailed context
func (l *Logger) ParseError(component, filePath string, err error, yamlContent string) {
	context := map[string]interface{}{
		"file":         filePath,
		"error":        err.Error(),
		"content_size": len(yamlContent),
	}
	
	// Add snippet of problematic content if available
	if len(yamlContent) > 0 {
		lines := strings.Split(yamlContent, "\n")
		maxLines := 10
		if len(lines) > maxLines {
			context["content_snippet"] = strings.Join(lines[:maxLines], "\n") + "..."
		} else {
			context["content_snippet"] = yamlContent
		}
	}

	l.Error(component, "YAML parsing failed", context)
}

// DiscoveryInfo logs discovery information
func (l *Logger) DiscoveryInfo(toolType, configPath, message string) {
	context := map[string]interface{}{
		"tool_type":   toolType,
		"config_path": configPath,
	}
	l.Info("Discovery", message, context)
}

// AnalysisError logs analysis errors with context
func (l *Logger) AnalysisError(toolType, configPath string, err error) {
	context := map[string]interface{}{
		"tool_type":   toolType,
		"config_path": configPath,
		"error":       err.Error(),
	}
	l.Error("Analysis", "Analysis failed", context)
}

// RecoveryAttempt logs recovery attempts during parsing
func (l *Logger) RecoveryAttempt(component, filePath, strategy string, success bool) {
	context := map[string]interface{}{
		"file":     filePath,
		"strategy": strategy,
		"success":  success,
	}
	
	if success {
		l.Info(component, "Recovery parsing succeeded", context)
	} else {
		l.Warn(component, "Recovery parsing failed", context)
	}
}

// Convenience functions for common use cases

// LogError is a convenience function for error logging
func LogError(component, message string, context map[string]interface{}) {
	GetLogger().Error(component, message, context)
}

// LogWarn is a convenience function for warning logging
func LogWarn(component, message string, context map[string]interface{}) {
	GetLogger().Warn(component, message, context)
}

// LogInfo is a convenience function for info logging
func LogInfo(component, message string, context map[string]interface{}) {
	GetLogger().Info(component, message, context)
}

// LogDebug is a convenience function for debug logging
func LogDebug(component, message string, context map[string]interface{}) {
	GetLogger().Debug(component, message, context)
}