package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

// HTTPErrorHandler more verbose version of the default error handler
func HTTPErrorHandler(err error, c echo.Context) {

	code := http.StatusInternalServerError
	msg := echo.Map{"message": http.StatusText(code)}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if he.Message != nil {
			msg["message"] = he.Message
		}
		if he.Internal != nil {
			msg["detail"] = he.Internal.Error() // check if stack trace available
		}
	}

	if code == http.StatusInternalServerError {
		printStackTrace(err)
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
	}
}

func printStackTrace(err error) {
	log.Errorf("%v", err)
	pErr := err
	if he, ok := err.(*echo.HTTPError); ok {
		pErr = he.Internal
	}
	fmt.Printf("%+v", pErr)
}
