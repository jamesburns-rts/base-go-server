package example

import (
	"github.com/jamesburns-rts/base-go-server/internal/call"
	"github.com/jamesburns-rts/base-go-server/internal/clients"
	"github.com/jamesburns-rts/base-go-server/internal/clients/ctxhttp"
	"github.com/jamesburns-rts/base-go-server/internal/config"
	"net/http"
)

// Client is used to make HTTP calls to the another service
type Client struct {
	HttpClient *http.Client
	BaseUrl    string
}

func NewClient(props config.ExampleClient) *Client {
	httpClient := props.HttpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if props.DebugLevel == "DEBUG" {
		httpClient = clients.AddDebugTransport(httpClient)
	} else {
		httpClient = clients.AddInfoTransport(httpClient)
	}

	return &Client{
		HttpClient: httpClient,
		BaseUrl:    props.Url,
	}
}

func (c *Client) GetExample(ctx call.Context, name string) (*Pokemon, error) {
	uri := c.BaseUrl + "/api/v2/pokemon/" + name

	// todo validate name since it is in url

	var pokemon Pokemon
	res, err := ctxhttp.GetJSON(ctx, c.HttpClient, uri, &pokemon)
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pokemon, nil
}
