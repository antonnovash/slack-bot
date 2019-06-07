package errors

import (
	"fmt"
)

// WrapError wraps error into one, that shows in what function it occurred, message and error value.
func WrapError(funcName string, message string, err error) error {
	return fmt.Errorf("[%s] %s: %v", funcName, message, err)
}
