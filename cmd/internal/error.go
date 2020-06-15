package internal

import "fmt"

type Error struct {
	Message string
	StatusCode int
}

func NewError(message string, statusCode int) *Error {
	return &Error{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s", e.Message)
}
