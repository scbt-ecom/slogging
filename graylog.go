package slogging

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/Graylog2/go-gelf/gelf"
	slogcommon "github.com/samber/slog-common"
)

type Converter func(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) (extra map[string]any)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// connection to graylog
	Writer *gelf.Writer

	// optional: customize json payload builder
	Converter Converter
	// optional: fetch attributes from context
	AttrFromContext []func(ctx context.Context) []slog.Attr

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	// internal
	hostname string
}

func (o Option) NewGraylogHandler() slog.Handler {
	if o.Level == nil {
		o.Level = LevelDebug
	}

	if o.AttrFromContext == nil {
		o.AttrFromContext = []func(ctx context.Context) []slog.Attr{}
	}

	if hostname, err := os.Hostname(); err == nil {
		o.hostname = hostname
	}

	return &GraylogHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*GraylogHandler)(nil)

type GraylogHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *GraylogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
	//return level >= h.option.Level.Level() || level == LevelFatal
}

func (h *GraylogHandler) Handle(ctx context.Context, record slog.Record) error {
	fromContext := slogcommon.ContextExtractor(ctx, h.option.AttrFromContext)
	extra := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, append(h.attrs, fromContext...), h.groups, &record)

	msg := &gelf.Message{
		Version:  "1.1",
		Host:     h.option.hostname,
		Short:    short(&record),
		TimeUnix: float64(record.Time.Unix()),
		Level:    LogLevels[record.Level],
		Extra:    extra,
	}

	// non-blocking
	go func() {
		_ = h.option.Writer.WriteMessage(msg)
	}()

	return nil
}

func (h *GraylogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GraylogHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *GraylogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return &GraylogHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

func short(record *slog.Record) string {
	msg := strings.TrimSpace(record.Message)
	if i := strings.IndexRune(msg, '\n'); i > 0 {
		return msg[:i]
	}

	return msg
}

const (
	XB3TraceID = "X-B3-TraceId"
)

var LogLevels = map[slog.Level]int32{
	slog.LevelDebug: 7,
	slog.LevelInfo:  6,
	slog.LevelWarn:  4,
	slog.LevelError: 3,
	LevelFatal:      2,
}
