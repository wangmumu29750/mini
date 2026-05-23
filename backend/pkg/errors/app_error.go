package errors

import "fmt"

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func Wrap(err error, status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Err: err}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
