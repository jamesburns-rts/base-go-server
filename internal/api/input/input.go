package input

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

// Bind reads the call's context body into i and validates required members
// i must be a pointer to a struct with 'json' tags for the desired fields - or xml I guess
// see util/input.go for creation of the validator
func Bind(c echo.Context, i interface{}) error {

	if err := c.Bind(i); err != nil {
		return errors.Wrap(err, "error while binding body")
	}
	if err := c.Validate(i); err != nil {
		return errors.Wrap(err, "invalid body for given struct")
	}
	return nil
}

// BindParams reads the call's context params and query params into i and validates the required members
// i must be a pointer to a struct with either 'param' or 'query' tags for the desired fields
// see util/input.go for creation of the validator
func BindParams(c echo.Context, i interface{}) error {
	if err := bindParams(c, i); err != nil {
		return errors.Wrap(err, "error while binding params")
	}
	if err := c.Validate(i); err != nil {
		return errors.Wrap(err, "invalid params for given struct")
	}
	return nil
}

// MustBeFloat should only be used after being validated by validator
func MustBeFloat(str string) float64 {
	fl, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Error(err, "Number that was assumed float is not")
	}
	return fl
}

func QueryParamList(c echo.Context, key string) (list []string) {
	params, ok := c.QueryParams()[key]
	if !ok {
		return list
	}
	for _, p := range params {
		for _, v := range strings.Split(p, ",") {
			if v != "" {
				list = append(list, v)
			}
		}
	}
	return list
}

// using reflection to actually read the context into the struct
func bindParams(c echo.Context, i interface{}) error {

	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Pointer not given")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("Pointer to struct not given")
	}

	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		valueField := rv.Field(i)
		switch valueField.Kind() {
		case reflect.Struct:
			if !valueField.Addr().CanInterface() {
				continue
			}

			iface := valueField.Addr().Interface()
			err := bindParams(c, iface)
			if err != nil {
				return errors.Wrap(err, "unable set sub struct in binding params")
			}
		}

		var paramValue string

		typeField := t.Field(i)
		if typeField.Tag.Get("query") != "" {
			paramValue = c.QueryParam(typeField.Tag.Get("query"))

		} else if typeField.Tag.Get("param") != "" {
			paramValue = c.Param(typeField.Tag.Get("param"))
		}

		if paramValue != "" {
			err := set(typeField.Type, valueField, paramValue)
			if err != nil {
				return errors.Wrap(err, "unable set value in binding params")
			}
		}
	}
	return nil
}

// use reflection to set the value based on the type
func set(t reflect.Type, f reflect.Value, value string) error {
	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := set(t.Elem(), ptr.Elem(), value)
		if err != nil {
			return errors.Wrap(err, "unable set value in binding param")
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return errors.Wrap(err, "unable to parse bool in binding param")
		}
		f.SetBool(v)
	case reflect.Int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return errors.Wrap(err, "unable to parse int in binding param")
		}
		f.SetInt(int64(v))
	default:
		return errors.New("unsupported type")
	}

	return nil
}
