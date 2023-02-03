package util

import (
	"github.com/labstack/gommon/log"
	"io"
	"runtime/debug"
)

// SafeClose after calling Close, it logs the error if one exists
// this should be used with defer statements
//
//	defer SafeClose(myObjectThatImplementsClose)
func SafeClose(c io.Closer) {
	if c == nil {
		return
	}
	err := c.Close()
	if err != nil {
		log.Error("error closing closer", err)
		debug.PrintStack()
	}
}
