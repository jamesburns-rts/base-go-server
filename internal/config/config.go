package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Netflix/go-env"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type (
	Application struct {
		Port          int    `env:"APPLICATION_PORT"`
		LocalHost     string `env:"APPLICATION_LOCAL_HOST"`
		LogLevel      string `env:"APPLICATION_LOG_LEVEL"`
		AuthEnabled   bool   `env:"APPLICATION_AUTH_ENABLED"`
		ExampleClient ExampleClient
	}

	ExampleClient struct {
		Url     string `env:"EXAMPLE_CLIENT_URL"`
		Timeout string `env:"EXAMPLE_CLIENT_TIMEOUT"`
	}
)

var Defaults = Application{
	Port:        8080,
	LocalHost:   "0.0.0.0",
	LogLevel:    "INFO",
	AuthEnabled: false,
	ExampleClient: ExampleClient{
		Url:     "http://client:8080/api",
		Timeout: "5s",
	},
}

// ReadProperties Reads the properties from the environment variables
func ReadProperties() (app Application, err error) {

	app = Defaults

	_, err = env.UnmarshalFromEnviron(&app)
	if err != nil {
		return app, err
	}

	// check some string values
	if _, err := time.ParseDuration(app.ExampleClient.Timeout); err != nil {
		return app, errors.Wrap(err, "invalid example client timeout given")
	}

	return app, err
}

func (app Application) EchoLogLevel() log.Lvl {
	switch strings.ToUpper(app.LogLevel) {
	default:
		return log.INFO
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OFF":
		return log.OFF
	}
}

func PrintProperties(props Application) {

	fmt.Println("\nApplication:")
	if err := printPropertiesRecursive(&props, "  "); err != nil {
		log.Warn("Unable to print rest of properties ", err)
	}
	fmt.Println()
}

func printPropertiesRecursive(v interface{}, indent string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Pointer not given")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("Pointer to struct not given")
	}

	rt := reflect.Indirect(rv).Type()

	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		valueField := rv.Field(i)
		switch valueField.Kind() {
		case reflect.Struct:
			if !valueField.Addr().CanInterface() {
				continue
			}

			iface := valueField.Addr().Interface()
			fmt.Println(indent + rt.Field(i).Name + ":")
			err := printPropertiesRecursive(iface, indent+"  ")
			if err != nil {
				return err
			}
		}

		typeField := t.Field(i)
		tag := typeField.Tag.Get("env")
		if tag == "" {
			continue
		}

		var value interface{}
		if typeField.Tag.Get("hide") != "" {
			value = "********"
		} else {
			value = valueField.Interface()
		}

		fmt.Println(indent+tag+":", value)
	}

	return nil
}
