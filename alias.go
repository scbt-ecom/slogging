package slogging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
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

func ResponseAttr(r *http.Response, start time.Time) []any {
	if r == nil {
		slog.Error("response is nil")
		return []any{}
	}

	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		r.Header.Set("Authorization", checkHeaderAuth(authHeader))
	}
	headers, _ := json.Marshal(r.Header)

	return getReqAttrsAsAny([]Attr{
		slog.String("url", r.Request.URL.String()),
		slog.String("method", r.Request.Method),
		slog.Int("statusCode", r.StatusCode),
		slog.String("headers", string(headers)),
		slog.String("body", string(body)),
		slog.Duration("duration", time.Since(start)),
	})
}

func RequestAttr(r *http.Request) []any {
	if r == nil {
		slog.Error("request is nil")
		return []any{}
	}

	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		r.Header.Set("Authorization", checkHeaderAuth(authHeader))
	}
	headers, _ := json.Marshal(r.Header)

	return getReqAttrsAsAny([]Attr{
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()),
		slog.String("headers", string(headers)),
		slog.String("body", string(body)),
	})
}

func checkHeaderAuth(header string) string {
	split := strings.Split(header, " ")
	if len(split) != 2 {
		return header
	}

	authData := split[1]
	if len(authData) >= 2 {
		authData = authData[len(authData)-2:]
	}

	return fmt.Sprintf("%s ***%s", split[0], authData)
}

func getReqAttrsAsAny(reqAttrs []Attr) []any {
	args := make([]any, len(reqAttrs))
	for i, attr := range reqAttrs {
		args[i] = attr
	}
	return args
}
