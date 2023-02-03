package ctxhttp

import "fmt"

type ResponseError struct {
	Code int
	Body string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Body)
}
