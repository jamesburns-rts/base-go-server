package api

import (
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/call"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/jamesburns-rts/base-go-server/internal/jwt"
	"github.com/jamesburns-rts/base-go-server/internal/log"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"strings"
	"time"
)

const UserContextKey = "user-context-key"

type ContextCreator func(e echo.Context) call.Context

func NewEchoContextCreator(props config.Application, dbc *sqlx.DB) ContextCreator {
	return func(e echo.Context) call.Context {
		ctx := &call.DefaultContext{}

		// the following is from controller constructor
		ctx.DatabaseConn = dbc

		// the following is based on the http request call being made
		ctx.Context = e.Request().Context()
		ctx.ContextID = e.Response().Header().Get(echo.HeaderXRequestID)
		ctx.CallerUserID = 0 // todo
		ctx.Logger = log.NewLogger(props.LogLevel, "ctx")
		ctx.Logger.AddAttributes(slog.Group("http",
			slog.String("method", e.Request().Method),
			slog.String("path", e.Request().URL.Path),
			slog.String("remoteAddress", e.Request().RemoteAddr),
			slog.String(echo.HeaderXForwardedFor, e.Request().Header.Get(echo.HeaderXForwardedFor)),
		))

		return ctx
	}

}

// Middleware all the configured middleware
// note the order is important for some of them
func Middleware(props config.Application) ([]echo.MiddlewareFunc, error) {
	middlewareJWT, err := createJWTMiddleware(props)
	if err != nil {
		return nil, err
	}

	return []echo.MiddlewareFunc{

		// if panic happens, don't crash the server
		middleware.Recover(),

		// add/get X-Request-ID header to identify each request
		middleware.RequestID(),

		// log each incoming request and response
		// note that this middleware consumes error and does not return it
		middlewareRequestLogger(props),

		// allow for gzipped requests
		middleware.Gzip(),

		// return 401 unless valid token exists
		// or route is public
		middlewareJWT,

		// only allow certain origins
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: strings.Split(props.CORSOrigins, ","),
		}),

		// default to json accept and content-types
		middlewareDefaultJSONContent,
	}, nil
}

func middlewareDefaultJSONContent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		header := e.Request().Header
		if header.Get(echo.HeaderContentType) == "" {
			header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		}
		if header.Get(echo.HeaderAccept) == "" {
			header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
		}
		return next(e)
	}
}

// middlewareRequestLogger is a custom version of echo.Logger which uses our structured logger
// this should be placed in the middleware stack early on - after RequestID but definitely before authentication
func middlewareRequestLogger(props config.Application) echo.MiddlewareFunc {
	contextCreator := NewEchoContextCreator(props, nil)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			ctx := contextCreator(e)

			// gather some attributes for the incoming request
			bytes := e.Request().Header.Get(echo.HeaderContentLength)
			if bytes == "" {
				bytes = "0"
			}

			ctx.Log().Info("Incoming Request",
				"bytes", bytes,
			)

			// record start time to calculate duration (ms)
			startTime := time.Now()

			// execute the next middleware (controller request handler is in there somewhere)
			err := next(e)
			if err != nil {
				// call echo context Error to set response properly
				e.Error(err)
			}

			duration := time.Since(startTime)
			res := e.Response()

			// the default context includes request id, route, method, etc
			attrs := []any{
				"status", res.Status,
				"ms", duration.Milliseconds(),
				"bytes", res.Size,
			}
			if err != nil {
				attrs = append(attrs, "err", err)
			}

			ctx.Log().Info("Request Status", attrs...)

			return nil // any err is consumed by this middleware
		}
	}
}

func createJWTMiddleware(props config.Application) (echo.MiddlewareFunc, error) {
	if !props.AuthEnabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc { return next }, nil
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(props.JWT.RASPublicKey))
	if err != nil {
		return nil, fmt.Errorf("invalid public key provided: %w", err)
	}

	leeway, err := time.ParseDuration(props.JWT.ExpirationLeeway)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt leeway provided: %w", err)
	}

	return echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if c.Request().Method == http.MethodOptions {
				return true
			}
			if strings.HasPrefix(c.Request().URL.Path, "/api/public/") {
				return true
			}
			return false
		},
		ContextKey: UserContextKey,
		ParseTokenFunc: func(_ echo.Context, auth string) (any, error) {
			return jwt.ParseAccessToken(auth, publicKey, leeway)
		},
	}), nil
}
