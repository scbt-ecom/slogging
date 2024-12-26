package slogging

import (
	"context"
	"github.com/Graylog2/go-gelf/gelf"
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
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
	defaultLevel      = LevelInfo
	defaultWithSource = true
	defaultSetDefault = true
)

// NewLogger opts can be
// InGraylog()
// SetLevel()
// WithSource()
// SetDefault()
func NewLogger(opts ...LoggerOption) *Logger {

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
	}

	stdHandler = NewTextHandler(os.Stdout, handlerOpts)

	if cfg.InGraylog == nil {
		l = New(stdHandler)
	} else {
		graylogHandler := Option{
			Level:     cfg.InGraylog.level,
			Writer:    cfg.InGraylog.w,
			Converter: DefaultConverter,
		}.NewGraylogHandler()

		graylogHandler = graylogHandler.WithAttrs([]Attr{
			slog.String("container_name", cfg.InGraylog.containerName)},
		)

		l = New(slogmulti.Fanout(stdHandler, graylogHandler))
	}

	if cfg.SetDefault {
		slog.SetDefault(l)
	}

	return l
}

type LoggerOption func(*LoggerConfig)

func WithAttrs(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}
