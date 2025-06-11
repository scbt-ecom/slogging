package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/scbt-ecom/slogging"
	amqpmw "github.com/scbt-ecom/slogging/amqp"
	ginmw "github.com/scbt-ecom/slogging/http/gin"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	log := slogging.NewLogger(
		slogging.SetLevel("debug"),
		slogging.InGraylog("localhost:12201", "debug", "application_name"),
		slogging.WithSource(true),
		slogging.SetDefault(true),
	)

	traceMW := ginmw.TraceMiddleware(log)

	amqpTraceMW := amqpmw.TraceMiddleware(log)

	r := gin.New()
	r.Use(traceMW)

	r.GET("/hello", ginHelloWorld)

	r.Run(":8080")
}

type TestStruct struct {
	A string
	B int
	C bool
}

func ginHelloWorld(c *gin.Context) {
	ctx := c.Request.Context()
	slogging.L(ctx).Info("hello world")
	slogging.L(ctx).Error("bye bye world")

	slog.Info("so good")
	slogging.L(context.Background()).Error("empty context test")

	//slogging.L(ctx).Fatal("fatal test = )")
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
