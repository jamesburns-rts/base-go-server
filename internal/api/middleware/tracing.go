package middleware

import (
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"strconv"
)

type pTracer struct{}

func SetupTracer(_ config.Application) *pTracer {

	t := opentracer.New(tracer.WithAgentAddr("host:port"))
	opentracing.SetGlobalTracer(t)

	return &pTracer{}
}

func (t *pTracer) Close() error {
	tracer.Stop()
	return nil
}

func Tracer(_ *pTracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// create a "root" Span from the given Context
			// and finish it when the request ends
			span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "echo.request")
			defer span.Finish()

			c.SetRequest(c.Request().WithContext(ctx))

			// propagate the trace in the Gin Context and process the request
			err := next(c)

			// add useful tags to your Trace
			span.SetTag("http.method", c.Request().Method)
			span.SetTag("http.status_code", strconv.Itoa(c.Response().Status))
			span.SetTag("http.url", c.Request().URL.Path)

			return err
		}
	}
}
