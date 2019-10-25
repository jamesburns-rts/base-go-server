package sample

import (
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/api/middleware"
	"github.com/jamesburns-rts/base-go-server/internal/api/response"
	. "github.com/jamesburns-rts/base-go-server/internal/model"
	"github.com/jamesburns-rts/base-go-server/internal/services/sample"

	"github.com/labstack/echo/v4"
)

// AddRoutes adds routes related to sample
func AddRoutes(g *echo.Group) {

	addSampleRoute(g.GET, "", getSample)
}

// addSampleRoute internal use to add sample context to route
func addSampleRoute(method middleware.MethodFunc, route string, controllerFunc sampleControllerFunc, middleware ...echo.MiddlewareFunc) {
	method(route, sampleContext(controllerFunc), middleware...)
}

type (
	// sampleController internal use struct containing necessary services to fulfill routes
	sampleController struct {
		sampleService sample.Service
	}

	// sampleControllerFunc internal use to define argument of addSampleRoute
	sampleControllerFunc func(s *sampleController, c echo.Context) error
)

// sampleContext internal use for wrapping a sampleControllerFunc with a sample controller
// request context, i.e. the database and queue connections
func sampleContext(f sampleControllerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {

		// get dependencies, such as database, from context
		client := middleware.GetExampleClient(context)

		controller := &sampleController{
			sampleService: sample.NewService(client),
		}

		return f(controller, context)
	}
}

// getSample GET /api/public/sample
// returns the sample specified by the path
func getSample(s *sampleController, c echo.Context) (err error) {

	var dto *SampleDTO
	if dto, err = s.sampleService.GetSample(); err != nil {
		return response.Error("error getting sample", err)
	}

	if dto == nil {
		return response.NotFound(fmt.Sprintf("sample not found"))
	}

	return response.Ok(c, dto.ToVm())
}
