package slogging

import (
	"golang.org/x/exp/constraints"
	"log/slog"
	"reflect"
	"time"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type (
	Logger         = slog.Logger
	Level          = slog.Level
	Record         = slog.Record
	Handler        = slog.Handler
	Attr           = slog.Attr
	HandlerOptions = slog.HandlerOptions
)

var (
	New            = slog.New
	NewJSONHandler = slog.NewJSONHandler
	NewTextHandler = slog.NewTextHandler
)

func IntAttr[T constraints.Integer](key string, value T) Attr {
	return slog.Int(key, int(value))
}

func FloatAttr[T constraints.Float](key string, value T) Attr {
	return slog.Float64(key, float64(value))
}

func TimeAttr(key string, time time.Time) Attr {
	return slog.String(key, time.Format("2006-01-02 15:04:05"))
}

func ErrAttr(err error) Attr {
	return slog.String("error", err.Error())
}

func StringAttr(key string, value string) Attr {
	return slog.String(key, value)
}

func AnyAttr(key string, s interface{}) Attr {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr && !v.IsNil() {
		s = v.Elem().Interface()
	}

	return slog.Any(key, s)
}
