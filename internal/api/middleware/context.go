package middleware

import (
	"github.com/jamesburns-rts/base-go-server/internal/clients/example"
	"github.com/labstack/echo/v4"
)

const exampleContextKey = "exampleContextKey"

// GetExampleClient retrieves the client instance from the context
// if not found, it panics
func GetExampleClient(c echo.Context) example.Client {
	d, ok := c.Get(exampleContextKey).(example.Client)
	if !ok || d == nil {
		panic("unable to get example client from context") // should be impossible
	}
	return d
}

// Middleware this adds client connection thread per call
// This middleware should be added to the top level echo instance
func ExampleClient(q example.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(exampleContextKey, q)
			return next(c)
		}
	}
}

// methodFunc is shorthand for echo method functions, e.g. group.POST, group.GET, etc.
type MethodFunc func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
