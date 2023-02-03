package api

import (
	"github.com/jamesburns-rts/base-go-server/internal"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Management struct {
}

func NewManagementController(props config.Application, ctxCreator ContextCreator) *Management {
	return &Management{}
}

func (c *Management) AddRoutes(g *echo.Group) {
	g.GET("/api/public/healthz", func(e echo.Context) error {
		// todo add db check?
		return e.JSONBlob(http.StatusOK, []byte(`{"status":"OK"}`))
	})
	g.GET("/api/public/version", func(e echo.Context) error {
		return e.JSON(http.StatusOK, internal.Versions)
	})
}
