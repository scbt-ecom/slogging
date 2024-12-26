package main

import (
	"errors"
	"fmt"
	"github.com/scbt-ecom/slogging"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	log := slogging.NewLogger(
		slogging.SetLevel("debug"),
		slogging.InGraylog("graylog:12201", "debug", "application_name"),
		slogging.WithSource(true),
		slogging.SetDefault(true),
	)

	tracemw := slogging.HTTPTraceMiddleware(log)

	http.HandleFunc("/", tracemw(helloWorld))

	a := &TestStruct{
		A: "das",
		B: 231,
		C: false,
	}

	abc := TestStruct{
		A: "ddasas",
		B: 1,
		C: false,
	}

	var sad *TestStruct
	fmt.Println(sad)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Info("example log message",
				slogging.ErrAttr(errors.New("example error message")),
				slogging.StringAttr("hello", "world"),
				slogging.IntAttr("bye", 12),
				slogging.FloatAttr("bye", 14.88),
				slogging.TimeAttr("time", time.Now()),
				slogging.AnyAttr("asd", a),
				slogging.AnyAttr("asddas", &a),
				slogging.AnyAttr("abc", abc),
				slogging.AnyAttr("abcabc", &abc),
				slogging.AnyAttr("sad", sad),
				slogging.AnyAttr("sadsad", &sad),
			)
		}
	}()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Info("лол почему?")
	}

}

type TestStruct struct {
	A string
	B int
	C bool
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	log := slogging.L(r.Context())
	log.Info("Тест HTTP ручки, тут должен быть TRACE заголовок")

	req, err := http.NewRequest("POST", "google.com", nil)
	if err != nil {
		slog.Info("message",
			slogging.ErrAttr(err))
	}

	req = slogging.RequestWithTraceHeaders(r.Context(), req)
	log.Info("headers",
		slogging.StringAttr("xb-3trace", req.Header.Get("X-B3-TraceId")))

	w.Write([]byte("Hello world"))
}
