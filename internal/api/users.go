package api

import (
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"github.com/jamesburns-rts/base-go-server/internal/svc"
	"github.com/jamesburns-rts/base-go-server/internal/vm"
	"github.com/jamesburns-rts/base-go-server/internal/vm/page"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Users struct {
	CreateCallContext ContextCreator
	UserService       *svc.UserService
}

func NewUsersController(props config.Application, ctxCreator ContextCreator) *Users {
	return &Users{
		CreateCallContext: ctxCreator,
		UserService:       svc.NewSampleService(props),
	}
}

func (c *Users) AddRoutes(g *echo.Group) {

	g.GET("/api/public/authenticate", func(e echo.Context) error {
		ctx := c.CreateCallContext(e)

		var login vm.Login
		if err := e.Bind(&login); err != nil {
			return err
		}

		tokens, err := c.UserService.Authenticate(ctx, login)
		if err != nil {
			return err
		}

		return e.JSON(http.StatusOK, tokens)
	})

	g.GET("/api/users/:userId/examples", func(e echo.Context) error {
		ctx := c.CreateCallContext(e)
		userId, err := strconv.Atoi(e.Param("userId"))
		if err != nil {
			return e.NoContent(http.StatusUnprocessableEntity)
		}

		pageParams, err := page.RequestFromParams(e.QueryParam("size"), e.QueryParam("page"), e.QueryParam("sort"))
		if err != nil {
			return err
		}

		examples, err := c.UserService.GetUserExamplesPage(ctx, userId, pageParams)
		if err != nil {
			return err
		}

		return e.JSON(http.StatusOK, examples)
	})
}
