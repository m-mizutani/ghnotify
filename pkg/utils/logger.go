package utils

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/m-mizutani/clog"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/masq"
	"golang.org/x/exp/slog"

	"github.com/m-mizutani/ghnotify/pkg/controller/cmd/option"
)

var (
	Logger = slog.New(clog.New(
		clog.WithWriter(os.Stdout),
		clog.WithLevel(slog.LevelInfo)),
	)
	mutex            = &sync.Mutex{}
	currentLogFormat = option.LogFormatConsole
)

func RenewLogger(w io.Writer, level option.LogLevel, format option.LogFormat) {
	filter := masq.New(masq.WithTag("secret"))

	var newLogger *slog.Logger
	switch format.Format() {
	case option.LogFormatConsole:
		newLogger = slog.New(clog.New(
			clog.WithWriter(w),
			clog.WithLevel(level.Level()),
			clog.WithReplaceAttr(filter),
		))

	case option.LogFormatJSON:
		newLogger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level:       level.Level(),
			ReplaceAttr: filter,
		}))
	}

	mutex.Lock()
	defer mutex.Unlock()
	Logger = newLogger
	currentLogFormat = format.Format()
}

func ErrLog(err error) any {
	if err == nil {
		return nil
	}

	attrs := []any{
		slog.String("message", err.Error()),
	}

	if goErr := goerr.Unwrap(err); goErr != nil {
		var values []any
		for k, v := range goErr.Values() {
			values = append(values, slog.Any(k, v))
		}
		attrs = append(attrs, slog.Group("values", values...))

		var stacktrace any
		if currentLogFormat == option.LogFormatJSON {
			var traces []string
			for _, st := range goErr.StackTrace() {
				traces = append(traces, fmt.Sprintf("%+v", st))
			}
			stacktrace = traces
		} else {
			stacktrace = goErr.StackTrace()
		}

		attrs = append(attrs, slog.Any("stacktrace", stacktrace))
	}

	return slog.Group("error", attrs...)
}
