package example

import (
	"net/http"
	"time"

	"github.com/jamesburns-rts/base-go-server/internal/config"
)

// Client is used to make HTTP calls to the uaa service
type (
	Client interface {
		GetExample() (str string, err error)
	}

	client struct {
		client      *http.Client
		UrlEndpoint string
	}
)

func NewClient(props config.Application) Client {
	timeout, _ := time.ParseDuration(props.ExampleClient.Timeout)
	return &client{
		client:      &http.Client{Timeout: timeout},
		UrlEndpoint: props.ExampleClient.Url,
	}
}

// GetTokenURL Gets the token url from the UAA endpoint
func (c *client) GetExample() (str string, err error) {
	return "hello there", nil
}
