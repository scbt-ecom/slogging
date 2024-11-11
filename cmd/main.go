package main

import (
	"errors"
	"github.com/scbt-ecom/slogging"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	log := abc.NewLogger(
		abc.SetLevel("debug"),
		abc.InGraylog("graylog:12201", "debug", "application_name"),
		abc.SetDefault(true),
	)

	tracemw := abc.HTTPTraceMiddleware(log)

	http.HandleFunc("/", tracemw(helloWorld))

	slogging.L(ctx).Info("example log message",
		slogging.ErrAttr(errors.New("example error message")),
		slogging.StringAttr("hello", "world"),
		slogging.IntAttr("bye", 12),
		slogging.AnyAttr("data", object),
		slogging.FloatAttr("bye", 14.88),
		slogging.TimeAttr("timestamp", time.Now()),
	)

	ctx = slogging.ContextWithLogger(ctx, slog.Default())

	log := slogging.L(ctx).With(
		slogging.StringAttr("module", "keycloak"))

	ctx = slogging.ContextWithLogger(ctx, log)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Info("лол почему?")
	}
}

//
//type TestStruct struct {
//	A string
//	B int
//	C bool
//}
//
//func helloWorld(w http.ResponseWriter, r *http.Request) {
//	log := abc.L(r.Context())
//	log.Info("Тест HTTP ручки, тут должен быть TRACE заголовок")
//
//	req, err := http.NewRequest("POST", "google.com", nil)
//	if err != nil {
//		slog.Info("message",
//			abc.ErrAttr(err))
//	}
//
//	req = abc.RequestWithTraceHeaders(r.Context(), req)
//	log.Info("headers",
//		abc.StringAttr("xb-3trace", req.Header.Get("X-B3-TraceId")))
//
//	w.Write([]byte("Hello world"))
//}
