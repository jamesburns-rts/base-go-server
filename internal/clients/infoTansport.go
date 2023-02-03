package clients

import (
	"github.com/jamesburns-rts/base-go-server/internal/call"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

// copied out of google.golang.org/api@v0.56.0/examples/debug.go

// InfoTransport http client layer that just prints the url and status code
type InfoTransport struct {
	Transport http.RoundTripper
}

// AddInfoTransport adds an InfoTransport layer to client using the provided logger
func AddInfoTransport(client *http.Client) *http.Client {
	clientCopy := *client
	clientCopy.Transport = &InfoTransport{
		Transport: client.Transport,
	}
	return &clientCopy
}

// RoundTrip implementation of http.RoundTripper
func (t *InfoTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	callStartTime := time.Now()
	res, err := transport.RoundTrip(req)
	status := 0
	if res != nil {
		status = res.StatusCode
	}
	callDuration := time.Since(callStartTime)

	logger := call.LoggerFromContext(req.Context())
	logger.Info("External API Call", slog.Group("req",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.Int("status", status),
		slog.Int64("ms", callDuration.Milliseconds()),
	))
	return res, err
}
