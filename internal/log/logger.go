package log

import (
	"errors"
	"golang.org/x/exp/slog"
	"os"
	"strings"
)

// Logger
// adding the AddAttributes method. The slog implementation has no
// way of adding attributes to an existing logger, so this adds that.
type Logger struct {
	*slog.Logger
}

func NewLogger(logLevel, source string) *Logger {
	var lvl slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "INFO":
		lvl = slog.LevelInfo
	default:
		lvl = slog.LevelInfo
	}
	handlerOptions := slog.HandlerOptions{Level: lvl}

	// could determine handler here (like use logfmt instead)
	h := slog.Handler(handlerOptions.NewJSONHandler(os.Stdout))
	h = h.WithAttrs([]slog.Attr{slog.String("logger", source)})

	// wrap logger so that attributes can be added on the fly
	return &Logger{Logger: slog.New(h)}
}

func (a *Logger) AddAttributes(args ...any) {
	a.Logger = a.Logger.With(args...)
}

// DefaultLogger is accessible through the exposed functions below it.
// It should avoid being used during a call context, but is necessary for some places.
var DefaultLogger *Logger

func Debug(msg string, i ...any)            { DefaultLogger.Debug(msg, i...) }
func Info(msg string, i ...any)             { DefaultLogger.Info(msg, i...) }
func Warn(msg string, i ...any)             { DefaultLogger.Warn(msg, i...) }
func Error(msg string, err error, i ...any) { DefaultLogger.Error(msg, err, i...) }
func Enabled(level slog.Level) bool         { return DefaultLogger.Enabled(level) }
func AddAttributes(attrs ...any)            { DefaultLogger.AddAttributes(attrs...) }

// Fatal special global function for global logger
// logs the message as an Error and terminates the program
func Fatal(msg string, i ...any) {
	err := errors.New(msg)
	if len(i) > 0 {
		iErr, ok := i[0].(error)
		if ok && iErr != nil {
			err = iErr
			i = i[1:]
		}
	}
	Error(msg, err, i...)
	os.Exit(1)
}

func init() {
	DefaultLogger = NewLogger("INFO", "default")
}
