package perrors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

type PluginError struct {
	errorMessage string
	err          error
	frame        xerrors.Frame
}

type stackError struct {
	*PluginError
}

func newError(cause error, errorMessage string) error {
	return &PluginError{
		errorMessage: errorMessage,
		err:          cause,
		frame:        xerrors.Caller(2),
	}
}

func Newf(format string, a ...interface{}) error {
	return newError(nil, fmt.Sprintf(format, a...))
}

func Wrapf(cause error, format string, a ...interface{}) error {
	return newError(cause, fmt.Sprintf(format, a...))
}

func (e *PluginError) Unwrap() error {
	return e.err
}

func Stack(err error) error {
	return &stackError{
		PluginError: &PluginError{
			errorMessage: err.Error(),
			err:          err,
			frame:        xerrors.Caller(1),
		},
	}
}

type printer struct {
	stacktrace *strings.Builder
}

func (p *printer) Print(args ...interface{}) {
	fmt.Fprint(p.stacktrace, args...)
}

func (p *printer) Printf(format string, args ...interface{}) {
	fmt.Fprintf(p.stacktrace, format, args...)
}

func (p *printer) Detail() bool {
	return true
}

func (e *PluginError) Error() string {
	p := &printer{
		stacktrace: &strings.Builder{},
	}
	p.stacktrace.WriteString(e.errorMessage + "\n - ")

	e.frame.Format(p)

	return p.stacktrace.String()
}
