package input

import (
	"github.com/jamesburns-rts/base-go-server/internal/util/page"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func PageFromParams(c echo.Context) (page.Request, error) {

	p := page.Request{
		PageSize:      20,
		PageNumber:    0,
		SortProperty:  "",
		SortDirection: "asc",
	}

	sizeStr := c.QueryParam("size")
	if sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return p, errors.Wrap(err, "invalid 'size' given, must be an integer")
		}
		if size <= 0 {
			return p, errors.New("'size' must be greater than zero")
		}
		p.PageSize = size
	}

	pageStr := c.QueryParam("page")
	if pageStr != "" {
		number, err := strconv.Atoi(pageStr)
		if err != nil {
			return p, errors.Wrap(err, "invalid 'page' given, must be an integer")
		}
		if number < 0 {
			return p, errors.New("'page' must be greater than or equal to zero ")
		}
		p.PageNumber = number
	}

	p.SortProperty = c.QueryParam("sortProperty")
	sortDirection := c.QueryParam("sortDirection")

	if sortDirection != "" && p.SortProperty == "" {
		return p, errors.New("must provide sortProperty when sortDirection is given")
	}

	switch strings.ToLower(sortDirection) {
	case "asc", "":
		p.SortDirection = "asc"
	case "desc":
		p.SortDirection = "desc"
	default:
		return p, errors.New("invalid sort direction given, must be either 'asc', or 'desc'")
	}

	return p, nil
}
