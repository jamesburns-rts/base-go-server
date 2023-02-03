// automatically validate structs based on tags

package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// docs list: https://godoc.org/gopkg.in/go-playground/validator.v9#hdr-Required
// code list: https://github.com/go-playground/validator/blob/v9/baked_in.go#L65
var modelValidator *validator.Validate

func init() {
	modelValidator = validator.New()

	// use json tags to get name
	modelValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate makes sure the underlying struct is valid
func Validate(i any) error {
	return modelValidator.Struct(i)
}
