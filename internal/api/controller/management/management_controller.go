package management

import (
	"github.com/dimiro1/health"
	"github.com/jamesburns-rts/base-go-server/internal/api/middleware"
	"github.com/jamesburns-rts/base-go-server/internal/api/response"
	"github.com/jamesburns-rts/base-go-server/internal/clients/example"
	versionUtil "github.com/jamesburns-rts/base-go-server/internal/util"
	"github.com/labstack/echo/v4"
)

func AddRoutes(g *echo.Group) {
	g.GET("/version", versionHandler)
	g.GET("/health", healthHandler)
}

func versionHandler(c echo.Context) error {
	return response.Ok(c, versionUtil.Versions)
}

// client health

type clientChecker struct {
	client example.Client
}

func (c *clientChecker) Check() health.Health {

	h := health.NewHealth()
	if _, err := c.client.GetExample(); err != nil {
		h.Up()
	} else {
		h.Down()
		h.AddInfo("error", err)
	}
	return h
}

func getClientChecker(client example.Client) health.Checker {
	return &clientChecker{client: client}
}

func healthHandler(c echo.Context) error {
	client := middleware.GetExampleClient(c)

	healthCheck := health.NewHandler()
	healthCheck.AddChecker("example-client", getClientChecker(client))
	healthCheck.ServeHTTP(c.Response(), c.Request())
	return nil
}
