package clients

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

// copied out of google.golang.org/api@v0.56.0/examples/debug.go

// DebugTransport http client layer that prints the url, status code, and request/response bodies
type DebugTransport struct {
	Transport http.RoundTripper
}

// AddDebugTransport adds an DebugTransport layer to client using the provided logger
func AddDebugTransport(client *http.Client) *http.Client {
	clientCopy := *client
	clientCopy.Transport = &DebugTransport{
		Transport: client.Transport,
	}
	return &clientCopy
}

// RoundTrip implementation of http.RoundTripper
func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var buf bytes.Buffer

	_, _ = os.Stdout.Write([]byte("\n[request]\n"))
	if req.Body != nil {
		req.Body = io.NopCloser(&readButCopy{req.Body, &buf})
	}
	_ = req.Write(os.Stdout)
	if req.Body != nil {
		req.Body = io.NopCloser(&buf)
	}
	_, _ = os.Stdout.Write([]byte("\n[/request]\n"))

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	res, err := transport.RoundTrip(req)

	fmt.Printf("[response]\n")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	} else {
		body := res.Body
		res.Body = nil
		_ = res.Write(os.Stdout)
		if body != nil {
			res.Body = io.NopCloser(&echoAsRead{body})
		}
	}

	return res, err
}

type echoAsRead struct {
	src io.Reader
}

// Read io.Reader implementation
func (r *echoAsRead) Read(p []byte) (int, error) {
	n, err := r.src.Read(p)
	if n > 0 {
		_, _ = os.Stdout.Write(p[:n])
	}
	if err == io.EOF {
		fmt.Printf("\n[/response]\n")
	}
	return n, err
}

type readButCopy struct {
	src io.Reader
	dst io.Writer
}

// Read io.Reader implementation
func (r *readButCopy) Read(p []byte) (int, error) {
	n, err := r.src.Read(p)
	if n > 0 {
		_, _ = r.dst.Write(p[:n])
	}
	return n, err
}
