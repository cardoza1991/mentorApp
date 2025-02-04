package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web"
)

// Log levels
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

// Logger represents our custom logger
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	console  bool
	level    string
	filepath string
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// InitLogger initializes the default logger
func InitLogger() error {
	var err error
	once.Do(func() {
		logPath := web.AppConfig.DefaultString("log::filepath", "logs/app.log")
		logLevel := web.AppConfig.DefaultString("log::level", "INFO")
		console := web.AppConfig.DefaultBool("log::console", true)

		// Ensure log directory exists
		logDir := filepath.Dir(logPath)
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		// Open log file
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		defaultLogger = &Logger{
			file:     file,
			console:  console,
			level:    logLevel,
			filepath: logPath,
		}
	})

	return err
}

// write writes a log entry to the configured outputs
func (l *Logger) write(entry *LogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Convert entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// Write to file
	if l.file != nil {
		if _, err := l.file.Write(append(jsonData, '\n')); err != nil {
			return err
		}
	}

	// Write to console if enabled
	if l.console {
		fmt.Printf("%s [%s] %s", entry.Timestamp, entry.Level, entry.Message)
		if len(entry.Fields) > 0 {
			fmt.Printf(" %+v", entry.Fields)
		}
		fmt.Println()
	}

	return nil
}

// log creates and writes a log entry
func (l *Logger) log(level, message string, fields ...interface{}) error {
	// Get caller information
	_, file, line, _ := runtime.Caller(2)
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	// Create log entry
	entry := &LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Caller:    caller,
		Fields:    make(map[string]interface{}),
	}

	// Process fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				entry.Fields[key] = fields[i+1]
			}
		}
	}

	return l.write(entry)
}

// Rotate rotates the log file if it exists
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return nil
	}

	// Close current file
	if err := l.file.Close(); err != nil {
		return err
	}

	// Rename current file with timestamp
	timestamp := time.Now().Format("20060102150405")
	newPath := fmt.Sprintf("%s.%s", l.filepath, timestamp)
	if err := os.Rename(l.filepath, newPath); err != nil {
		return err
	}

	// Open new file
	file, err := os.OpenFile(l.filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.file = file
	return nil
}

// Public logging functions

func Debug(message string, fields ...interface{}) {
	if defaultLogger == nil {
		InitLogger()
	}
	defaultLogger.log(LevelDebug, message, fields...)
}

func Info(message string, fields ...interface{}) {
	if defaultLogger == nil {
		InitLogger()
	}
	defaultLogger.log(LevelInfo, message, fields...)
}

func Warn(message string, fields ...interface{}) {
	if defaultLogger == nil {
		InitLogger()
	}
	defaultLogger.log(LevelWarn, message, fields...)
}

func Error(message string, fields ...interface{}) {
	if defaultLogger == nil {
		InitLogger()
	}
	defaultLogger.log(LevelError, message, fields...)
}

func Fatal(message string, fields ...interface{}) {
	if defaultLogger == nil {
		InitLogger()
	}
	defaultLogger.log(LevelFatal, message, fields...)
	os.Exit(1)
}

// GetWriter returns an io.Writer interface for the logger
func GetWriter() io.Writer {
	if defaultLogger == nil {
		InitLogger()
	}
	return defaultLogger.file
}

// SetLevel sets the logging level
func SetLevel(level string) {
	if defaultLogger == nil {
		InitLogger()
	}
	defaultLogger.level = level
}

// Close closes the logger and any associated resources
func Close() error {
	if defaultLogger != nil && defaultLogger.file != nil {
		return defaultLogger.file.Close()
	}
	return nil
}
