package logger

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

const (
	LoggerTimeFormat = "02/01/06-15:04:05.000"
)

type Lgr interface {
	Debug(msg string)
	DebugF(msg string, v ...interface{})
	Info(msg string)
	InfoF(msg string, v ...interface{})
	Warn(msg string)
	WarnF(msg string, v ...interface{})
	Error(err error)
	ErrorC(err error, errCode int)
}

type Logger struct {
	// silent determines if logger will log to stdout
	silent bool
	caller bool
	// console-output determines the type of output for console
	consoleOutputType ConsoleOutputType
	loggers           []io.Writer
	zero              zerolog.Logger
	lvl               zerolog.Level
	extraField        []ExtraField
}

type ConsoleOutputType int

const (
	ConsoleOutputTypeJSON ConsoleOutputType = iota
	ConsoleOutputTypeText
)

// NewLogger returns initialized Logger object
func NewLogger() *Logger {
	l := new(Logger)
	l.silent = false
	l.caller = false
	l.lvl = zerolog.InfoLevel
	l.consoleOutputType = ConsoleOutputTypeJSON
	zerolog.TimestampFieldName = "t"
	zerolog.TimeFieldFormat = LoggerTimeFormat
	zerolog.DurationFieldUnit = time.Millisecond
	l.zero = l.getLogger()
	return l
}

func (l Logger) Child() *Logger {
	newLogger := l
	newLogger.zero = l.zero
	newLogger.extraField = []ExtraField{}
	return &newLogger
}

// WithType sets the type of the logger
func (l *Logger) WithType(t ConsoleOutputType) *Logger {
	l.consoleOutputType = t
	l.zero = l.getLogger()
	return l
}

// WithLevelDebug sets the logger to log debug lvl messages
func (l *Logger) WithLevelDebug() *Logger {
	l.lvl = zerolog.DebugLevel
	l.zero = l.getLogger()
	return l
}

// WithLevelInfo sets the logger to log info lvl messages
func (l *Logger) WithLevelInfo() *Logger {
	l.lvl = zerolog.InfoLevel
	l.zero = l.getLogger()
	return l
}

// WithLevelWarn sets the logger to log warn lvl messages
func (l *Logger) WithLevelWarn() *Logger {
	l.lvl = zerolog.WarnLevel
	l.zero = l.getLogger()
	return l
}

// WithLevelError sets the logger to log error lvl messages
func (l *Logger) WithLevelError() *Logger {
	l.lvl = zerolog.ErrorLevel
	l.zero = l.getLogger()
	return l
}

// WithLevel sets the logger to log messages of the specified level
func (l *Logger) WithLevel(level string) *Logger {
	switch level {
	case "debug":
		l.lvl = zerolog.DebugLevel
	case "info":
		l.lvl = zerolog.InfoLevel
	case "warn":
		l.lvl = zerolog.WarnLevel
	case "error":
		l.lvl = zerolog.ErrorLevel
	default:
		l.lvl = zerolog.InfoLevel
	}
	l.zero = l.getLogger()
	return l
}

// WithWriter appends writer for logging
func (l *Logger) WithWriter(writer io.Writer) *Logger {
	l.loggers = append(l.loggers, writer)
	l.zero = l.getLogger()
	return l
}

// WithFile opens or creates a file that logs will be written to
func (l *Logger) WithFile(path string) *Logger {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	l.loggers = append(l.loggers, f)
	l.zero = l.getLogger()
	return l
}

// WithStr adds string key to the logger
func (l *Logger) WithStr(key, val string) *Logger {
	l.zero = l.zero.With().Str(key, val).Logger()
	return l
}

// WithCaller adds caller key to the log
func (l *Logger) WithCaller() *Logger {
	l.zero = l.zero.With().CallerWithSkipFrameCount(4).Logger()
	return l
}

// WithSilent set silent mode for logger
func (l *Logger) WithSilent(value bool) *Logger {
	l.silent = value
	l.zero = l.getLogger()
	return l
}

func (l *Logger) Logger() zerolog.Logger {
	return l.zero
}

// Debug is wrapper around zerolog Debug method
func (l *Logger) Debug(msg string) {
	l.fields(l.zero.Debug()).Msg(msg)
}

// DebugF is wrapper around zerolog Debug method
func (l *Logger) DebugF(msg string, v ...interface{}) {
	l.fields(l.zero.Debug()).Msgf(msg, v...)
}

// Info is wrapper around zerolog Info method
func (l *Logger) Info(msg string) {
	l.fields(l.zero.Info()).Msg(msg)
}

// InfoF is wrapper around zerolog Info method
func (l *Logger) InfoF(msg string, v ...interface{}) {
	l.fields(l.zero.Info()).Msgf(msg, v...)
}

// Warn is wrapper around zerolog Warn method
func (l *Logger) Warn(msg string) {
	l.fields(l.zero.Warn()).Msg(msg)
}

// WarnF is wrapper around zerolog Warn method
func (l *Logger) WarnF(msg string, v ...interface{}) {
	l.fields(l.zero.Warn()).Msgf(msg, v...)
}

// Error is wrapper around zerolog Error method
func (l *Logger) Error(err error) {
	l.fields(l.zero.Error().Err(err)).Msg("")
}

// ErrorC is wrapper around zerolog Error method
func (l *Logger) ErrorC(err error, errCode int) {
	l.fields(l.zero.Error().Err(err).Int("c", errCode)).Msg("")
}

// getLogger constructs zerolog.Logger from Logger type
func (l *Logger) getLogger() zerolog.Logger {
	loggers := l.loggers
	if !l.silent {
		var stdoutLogger io.Writer

		switch l.consoleOutputType {
		case ConsoleOutputTypeText:
			stdoutLogger = zerolog.ConsoleWriter{Out: os.Stdout}
		case ConsoleOutputTypeJSON:
			stdoutLogger = os.Stdout
		default:
			panic("unknown console output type")
		}

		if stdoutLogger != nil {
			loggers = append(loggers, stdoutLogger)
		}
	}
	writer := io.MultiWriter(loggers...)
	logger := zerolog.New(writer).Level(l.lvl).With().Timestamp().Logger()

	return logger
}

type ExtraField interface {
	Field(e *zerolog.Event) *zerolog.Event
}

func (l *Logger) fields(e *zerolog.Event) *zerolog.Event {
	for _, f := range l.extraField {
		e = f.Field(e)
	}
	return e
}
