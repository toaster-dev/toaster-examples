package xerrors

import (
	"errors"
)

type ErrorKind error

var ErrNotFound ErrorKind = errors.New("not found")

type StructuredError struct {
	Kind         ErrorKind
	Message      string
	WrappedError error
}

func (err StructuredError) Error() string {
	return err.Message
}

func (err StructuredError) Is(target error) bool {
	return err.Kind == target
}

func (err StructuredError) Unwrap() error {
	return err.WrappedError
}

func Error(kind ErrorKind, message string) error {
	return StructuredError{
		Kind:    kind,
		Message: message,
	}
}

func WrapError(kind ErrorKind, message string, err error) error {
	return StructuredError{
		Kind:         kind,
		Message:      message,
		WrappedError: err,
	}
}
