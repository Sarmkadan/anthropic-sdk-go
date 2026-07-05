// Package logging provides structured logging support using Microsoft.Extensions.Logging
// with slog as the underlying implementation.
//
// This package wraps Microsoft.Extensions.Logging.ILogger to provide structured
// logging with named parameters for better observability and debugging.
package logging

import (
	"context"
	"log/slog"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

// Logger is a wrapper around Microsoft.Extensions.Logging.ILogger that provides
// structured logging capabilities.
type Logger struct {
	logr.Logger
	name string
}

// New creates a new Logger with the given name.
func New(name string) *Logger {
	base := stdr.New(slog.Default().Handler())
	return &Logger{
		Logger: base.WithName(name),
		name:   name,
	}
}

// WithName creates a new Logger with an additional name segment.
func (l *Logger) WithName(name string) *Logger {
	return &Logger{
		Logger: l.Logger.WithName(name),
		name:   l.name + "." + name,
	}
}

// LogLevel represents the severity level for log messages.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// FromSlogLevel converts slog.Level to LogLevel.
func FromSlogLevel(level slog.Level) LogLevel {
	switch {
	case level >= slog.LevelError:
		return LogLevelError
	case level >= slog.LevelWarn:
		return LogLevelWarning
	case level >= slog.LevelInfo:
		return LogLevelInfo
	default:
		return LogLevelDebug
	}
}

// ToSlogLevel converts LogLevel to slog.Level.
func ToSlogLevel(level LogLevel) slog.Level {
	switch level {
	case LogLevelError:
		return slog.LevelError
	case LogLevelWarning:
		return slog.LevelWarn
	case LogLevelInfo:
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}

// Log is a low-level logging method that accepts a log level and message.
func (l *Logger) Log(level LogLevel, msg string, keysAndValues ...any) {
	l.Logger.V(int(level)).Info(msg, keysAndValues...)
}

// LogDebug logs a debug message.
func (l *Logger) LogDebug(msg string, args ...any) {
	l.Logger.V(int(LogLevelDebug)).Info(msg, args...)
}

// LogInfo logs an informational message.
func (l *Logger) LogInfo(msg string, args ...any) {
	l.Logger.V(int(LogLevelInfo)).Info(msg, args...)
}

// LogWarning logs a warning message.
func (l *Logger) LogWarning(msg string, args ...any) {
	l.Logger.V(int(LogLevelWarning)).Info(msg, args...)
}

// LogError logs an error message.
func (l *Logger) LogError(msg string, err error, args ...any) {
	args = append(args, "error", err)
	l.Logger.Error(err, msg, args...)
}

// LogWithContext logs a message with context from the provided context.Context.
func (l *Logger) LogWithContext(ctx context.Context, level LogLevel, msg string, args ...any) {
	if l == nil {
		return
	}

	// Extract request ID or other context values if available
	if reqID := ctx.Value(struct{}{}); reqID != "" {
		args = append(args, "request_id", reqID)
	}

	switch level {
	case LogLevelDebug:
		l.LogDebug(msg, args...)
	case LogLevelInfo:
		l.LogInfo(msg, args...)
	case LogLevelWarning:
		l.LogWarning(msg, args...)
	case LogLevelError:
		l.Logger.Error(nil, msg, args...)
	}
}

// IsDebugEnabled returns true if debug logging is enabled.
func (l *Logger) IsDebugEnabled() bool {
	return l.Logger.V(int(LogLevelDebug)).Enabled()
}

// IsInfoEnabled returns true if info logging is enabled.
func (l *Logger) IsInfoEnabled() bool {
	return l.Logger.V(int(LogLevelInfo)).Enabled()
}

// IsWarningEnabled returns true if warning logging is enabled.
func (l *Logger) IsWarningEnabled() bool {
	return l.Logger.V(int(LogLevelWarning)).Enabled()
}

// IsErrorEnabled returns true if error logging is enabled.
func (l *Logger) IsErrorEnabled() bool {
	return true // Errors should always be logged
}
