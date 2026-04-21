package errors

import "errors"

type Code string

const (
	ErrCodeNotFound   Code = "NOT_FOUND"
	ErrCodeValidation Code = "VALIDATION"
	ErrCodeInternal   Code = "INTERNAL"
	ErrCodeExecution  Code = "EXECUTION"
)

type TONError struct {
	Code    Code
	Message string
	Err     error
}

func (e *TONError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *TONError) Unwrap() error {
	return e.Err
}

func New(code Code, message string, err error) *TONError {
	return &TONError{Code: code, Message: message, Err: err}
}

func NewNotFound(message string) *TONError {
	return New(ErrCodeNotFound, message, nil)
}

func NewValidation(message string) *TONError {
	return New(ErrCodeValidation, message, nil)
}

func NewInternal(err error) *TONError {
	return New(ErrCodeInternal, "Internal error", err)
}

func Wrap(err error, code Code, message string) *TONError {
	if err == nil {
		return nil
	}
	return New(code, message, err)
}

func NewExecution(message string) *TONError {
	return New(ErrCodeExecution, message, nil)
}

var ErrToolNotFound = errors.New("tool not found")