// common return statuses and messages
package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Ok returns a 200 status
func Ok(c echo.Context, body interface{}) error {
	err := c.JSON(http.StatusOK, body)
	if err != nil {
		return InternalError("unable to write to json", err)
	}
	return nil
}

// Created returns a 201 status
func Created(c echo.Context, body interface{}) error {
	err := c.JSON(http.StatusCreated, body)
	if err != nil {
		return InternalError("unable to write to json", err)
	}
	return nil
}

// NoContent returns a 204 with no body
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// BadRequest 400 makes nice bad request messages
func BadRequest(message string, errs ...error) error {
	err := echo.NewHTTPError(http.StatusBadRequest, message)
	if len(errs) > 0 {
		err.Internal = err
	}
	return err
}

// Unauthorized 401 makes nice forbidden messages
func Unauthorized() error {
	return echo.NewHTTPError(http.StatusUnauthorized)
}

// Forbidden 403 makes nice forbidden messages
func Forbidden() error {
	return echo.NewHTTPError(http.StatusForbidden)
}

// NotFound 404 makes nice not found messages
func NotFound(message string) error {
	err := echo.NewHTTPError(http.StatusNotFound, message)
	return err
}

// InternalError wrapper around a 500, prints a stack trace
func InternalError(message string, errs ...error) error {
	err := echo.NewHTTPError(http.StatusInternalServerError, message)
	if len(errs) > 0 {
		err.Internal = errs[0]
	}
	return err
}

// BadRequestError error type that should be returned as a bad request
type BadRequestError interface {
	BadRequestMessage() string
	Detail() error
}

// Error parses err and returns the proper code based on the available interfaces
func Error(context string, err error) error {

	switch v := err.(type) {
	case BadRequestError:
		return BadRequest(context+": "+v.BadRequestMessage(), v.Detail())
	default:
		return InternalError(context, err)
	}
}
