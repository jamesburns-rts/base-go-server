// Package ctxhttp provides helper functions for performing context-aware HTTP requests.
// modified version of "golang.org/x/net/context/ctxhttp"
package ctxhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jamesburns-rts/base-go-server/internal/util"
	"io"
	"net/http"
)

// Do sends an HTTP request with the provided http.Client and returns
// an HTTP response.
//
// If the client is nil, http.DefaultClient is used.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func Do(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req.WithContext(ctx))
	// If we got an error, and the context has been canceled,
	// the context's error is probably more useful.
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}
	return resp, err
}

// GetRaw issues a GET request
func GetRaw(ctx context.Context, client *http.Client, url string) (*http.Response, string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	res, err := Do(ctx, client, req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to make GET request: %w", err)
	}

	defer util.SafeClose(res.Body)
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return res, "", fmt.Errorf("failed to get the body %w", err)
	}

	return res, string(bodyBytes), nil
}

// GetJSONWithParams issues a GET request via the GetPage function and unmarshals the body
func GetJSONWithParams(ctx context.Context, client *http.Client, url string, params any, bodyDest any) (*http.Response, error) {
	url, err := AddQueryParameters(url, params)
	if err != nil {
		return nil, fmt.Errorf("adding query params: %w", err)
	}
	return GetJSON(ctx, client, url, bodyDest)
}

// GetJSON issues a GET request via the GetPage function and unmarshals the body
func GetJSON(ctx context.Context, client *http.Client, url string, bodyDest any) (*http.Response, error) {
	return doWithJSONBody(ctx, client, http.MethodGet, url, nil, bodyDest)
}

// Post issues a POST request via the Do function.
func Post(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	return doWithRequestBody(ctx, client, http.MethodPost, url, bodyType, body)
}

// PostJSON marshals the body into json and issues a POST request via the Post function
func PostJSON(ctx context.Context, client *http.Client, url string, body any, bodyDest any) (*http.Response, error) {
	return doWithJSONBody(ctx, client, http.MethodPost, url, body, bodyDest)
}

// Put issues a PUT request via the Do function.
func Put(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	return doWithRequestBody(ctx, client, http.MethodPut, url, bodyType, body)
}

// PutJSON marshals the body into json and issues a PUT request via the Put function
func PutJSON(ctx context.Context, client *http.Client, url string, body any, bodyDest any) (*http.Response, error) {
	return doWithJSONBody(ctx, client, http.MethodPut, url, body, bodyDest)
}

// Patch issues a PATCH request via the Do function.
func Patch(ctx context.Context, client *http.Client, url string, bodyType string, body io.Reader) (*http.Response, error) {
	return doWithRequestBody(ctx, client, http.MethodPatch, url, bodyType, body)
}

// PatchJSON marshals the body into json and issues a PATCH request via the Patch function
func PatchJSON(ctx context.Context, client *http.Client, url string, body any, bodyDest any) (*http.Response, error) {
	return doWithJSONBody(ctx, client, http.MethodPatch, url, body, bodyDest)
}

// Delete issues a DELETE request via the Do function.
func Delete(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	res, err := doWithRequestBody(ctx, client, http.MethodDelete, url, "", nil)
	if err == nil && (res.StatusCode < 200 || res.StatusCode > 299) {
		defer util.SafeClose(res.Body)
		responseBody, _ := io.ReadAll(res.Body)
		err = &ResponseError{Code: res.StatusCode, Body: string(responseBody)}
	}
	return res, err
}

func doWithRequestBody(ctx context.Context, client *http.Client, method string, url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("Accept", bodyType)
	return Do(ctx, client, req)
}

func doWithJSONBody(ctx context.Context, client *http.Client, method string, url string, body any, bodyDest any) (_ *http.Response, err error) {
	var b []byte
	if body == nil {
		b = nil
	} else if s, ok := body.(string); ok {
		b = []byte(s)
	} else {
		b, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	res, err := doWithRequestBody(ctx, client, method, url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("failed to make %s request: %w", method, err)
	}
	defer util.SafeClose(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		responseBody, _ := io.ReadAll(res.Body)
		return res, &ResponseError{Code: res.StatusCode, Body: string(responseBody)}
	}

	if bodyDest != nil {
		err = json.NewDecoder(res.Body).Decode(bodyDest)
		if err != nil {
			return res, fmt.Errorf("failed to unmarshal body: %w", err)
		}
	}
	return res, nil
}
