package util

import (
	"github.com/go-playground/validator"
)

// InputValidator used to validate input from user
type InputValidator struct {
	// docs list: https://godoc.org/gopkg.in/go-playground/validator.v9#hdr-Required
	// code list: https://github.com/go-playground/validator/blob/v9/baked_in.go#L65
	validator *validator.Validate
}

// Validator is the initialized validator
var Validator = &InputValidator{validator.New()}

// Validate makes sure the underlying struct is valid
func (cv *InputValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
