package slogging

import (
	"context"
	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
	"time"
)

type LoggerConfig struct {
	Level      Level
	WithSource bool
	IsJSON     bool
	SetDefault bool
	InGraylog  *gelfData
}

type gelfData struct {
	w             *gelf.Writer
	level         Level
	containerName string
}

const (
	defaultLevel      = LevelDebug
	defaultWithSource = true
	defaultSetDefault = true
)

// NewLogger opts can be
// InGraylog()
// SetLevel()
// WithSource()
// SetDefault()
func NewLogger(opts ...LoggerOption) *SLogger {

	cfg := &LoggerConfig{
		Level:      defaultLevel,
		WithSource: defaultWithSource,
		SetDefault: defaultSetDefault,
		InGraylog:  nil,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	var l *Logger

	var stdHandler Handler
	handlerOpts := &HandlerOptions{
		AddSource: cfg.WithSource,
		Level:     cfg.Level,
		//ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		//	if a.Key == "X-B3-Order" && len(groups) == 0 {
		//		return slog.Attr{} // Удаляем поле
		//	}
		//	return a
		//},
	}

	//stdHandler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	stdHandler = NewTextHandler(os.Stdout, handlerOpts)

	if cfg.InGraylog == nil {
		l = New(stdHandler)
	} else {
		sloggraylog.SourceKey = "reference"
		graylogHandler := Option{
			Level:     cfg.InGraylog.level,
			Writer:    cfg.InGraylog.w,
			Converter: sloggraylog.DefaultConverter,
			AddSource: cfg.WithSource,
		}.NewGraylogHandler()

		graylogHandler = graylogHandler.WithAttrs([]Attr{
			slog.String("container_name", cfg.InGraylog.containerName)},
		)

		l = New(slogmulti.Fanout(stdHandler, graylogHandler))
	}

	if cfg.SetDefault {
		slog.SetDefault(l)
	}

	return &SLogger{
		Logger: l,
	}
}

type LoggerOption func(*LoggerConfig)

type SLogger struct {
	*slog.Logger
}

func (l *SLogger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
	time.Sleep(1 * time.Second)
	os.Exit(1)
}
