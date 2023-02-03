package ctxhttp

import (
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// AddQueryParameters adds the parameters in params as URL query parameters to s. params
// must be a struct whose fields may contain "url" tags.
func AddQueryParameters(s string, params any) (string, error) {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(params)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
