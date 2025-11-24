package logger

import (
	"encoding/json"
	"os"
	"time"
)

// Level represents log level
type Level string

const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

// Stage represents work stage
type Stage string

const (
	StageLL   Stage = "LL"
	StagePRP  Stage = "PRP"
	StageTF   Stage = "TF"
	StagePM1  Stage = "P-1"
	StagePP1  Stage = "P+1"
	StageECM  Stage = "ECM"
	StageIdle Stage = "idle"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp string  `json:"timestamp"`
	Level     Level   `json:"level"`
	Worker    *int    `json:"worker,omitempty"`
	Stage     *Stage  `json:"stage,omitempty"`
	Exponent  *uint64 `json:"exponent,omitempty"`
	Iteration *uint64 `json:"iteration,omitempty"`
	Progress  *float64 `json:"progress,omitempty"`
	Message   string  `json:"message"`
}

// Logger provides structured JSON logging
type Logger struct {
	output *os.File
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		output: os.Stdout,
	}
}

// NewFileLogger creates a logger that writes to a file
func NewFileLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{
		output: file,
	}, nil
}

// Close closes the logger output
func (l *Logger) Close() error {
	if l.output != os.Stdout && l.output != os.Stderr {
		return l.output.Close()
	}
	return nil
}

// log writes a log entry
func (l *Logger) log(level Level, entry LogEntry) {
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	entry.Level = level

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	l.output.Write(data)
	l.output.WriteString("\n")
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.log(LevelDebug, LogEntry{Message: message})
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.log(LevelInfo, LogEntry{Message: message})
}

// Warning logs a warning message
func (l *Logger) Warning(message string) {
	l.log(LevelWarning, LogEntry{Message: message})
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.log(LevelError, LogEntry{Message: message})
}

// WorkerProgress logs worker progress
func (l *Logger) WorkerProgress(worker int, stage Stage, exponent uint64, iteration, total uint64, message string) {
	progress := float64(iteration) / float64(total)
	l.log(LevelInfo, LogEntry{
		Worker:    &worker,
		Stage:     &stage,
		Exponent:  &exponent,
		Iteration: &iteration,
		Progress:  &progress,
		Message:   message,
	})
}

// WorkerError logs a worker error
func (l *Logger) WorkerError(worker int, stage Stage, exponent uint64, message string) {
	l.log(LevelError, LogEntry{
		Worker:   &worker,
		Stage:    &stage,
		Exponent: &exponent,
		Message:  message,
	})
}

// PrimeNetLog logs PrimeNet communication
func (l *Logger) PrimeNetLog(level Level, message string) {
	l.log(level, LogEntry{Message: message})
}

// PerformanceLog logs performance metrics
func (l *Logger) PerformanceLog(worker int, stage Stage, exponent uint64, duration time.Duration, message string) {
	durationMs := float64(duration.Nanoseconds()) / 1e6
	l.log(LevelInfo, LogEntry{
		Worker:   &worker,
		Stage:    &stage,
		Exponent: &exponent,
		Progress: &durationMs,
		Message:  message,
	})
}

// Global logger instance
var defaultLogger = NewLogger()

// Debug logs a debug message using default logger
func Debug(message string) {
	defaultLogger.Debug(message)
}

// Info logs an info message using default logger
func Info(message string) {
	defaultLogger.Info(message)
}

// Warning logs a warning message using default logger
func Warning(message string) {
	defaultLogger.Warning(message)
}

// Error logs an error message using default logger
func Error(message string) {
	defaultLogger.Error(message)
}

// WorkerProgress logs worker progress using default logger
func WorkerProgress(worker int, stage Stage, exponent uint64, iteration, total uint64, message string) {
	defaultLogger.WorkerProgress(worker, stage, exponent, iteration, total, message)
}

// WorkerError logs a worker error using default logger
func WorkerError(worker int, stage Stage, exponent uint64, message string) {
	defaultLogger.WorkerError(worker, stage, exponent, message)
}

// PrimeNetLog logs PrimeNet communication using default logger
func PrimeNetLog(level Level, message string) {
	defaultLogger.PrimeNetLog(level, message)
}

// PerformanceLog logs performance metrics using default logger
func PerformanceLog(worker int, stage Stage, exponent uint64, duration time.Duration, message string) {
	defaultLogger.PerformanceLog(worker, stage, exponent, duration, message)
}

