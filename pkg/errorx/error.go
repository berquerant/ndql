package errorx

import (
	"encoding/json"
	"fmt"
)

type ExitError struct {
	err  error
	code int
}

func WithExitCode(err error, code int) *ExitError {
	return &ExitError{
		err:  err,
		code: code,
	}
}

func (err *ExitError) Error() string { return fmt.Sprintf("%v: exit code %d", err.err, err.code) }
func (err *ExitError) Unwrap() error { return err.err }
func (err *ExitError) Code() int     { return err.code }

// WithValue append a default format string of the value to the error.
func WithValue(err error, v any) error { return fmt.Errorf("%w: %v", err, v) }

// WithVerbose append a verbose format string of the value to the error.
func WithVerbose(err error, v any) error { return fmt.Errorf("%w: %#v", err, v) }

// WithJSON append a JSON string of the value to the error.
func WithJSON(err error, v any) error {
	b, _ := json.Marshal(v)
	return fmt.Errorf("%w: %s", err, b)
}

// AsString converts the error into string.
func AsString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
