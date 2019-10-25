package api

import (
	"github.com/jamesburns-rts/base-go-server/internal/api/controller/management"
	"github.com/jamesburns-rts/base-go-server/internal/api/controller/sample"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo) {
	//sample.AddRoutes(e.Group("/api/sample", middleware.IsAuthed))
	sample.AddRoutes(e.Group("/api/public/sample"))
	management.AddRoutes(e.Group("/management"))
}
