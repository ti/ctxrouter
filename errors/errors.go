package errors

import (
	"fmt"
)

//Error Default HTTP Error
type Error struct {
	Message string   `json:"error,omitempty"`
	Code    Code     `json:"code,omitempty"`
	Details []Detail `json:"details,omitempty"`
}

// New returns a Status representing c and msg.
func New(c Code, msg string) *Error {
	if c == OK {
		return nil
	}
	return &Error{
		Message: msg,
		Code:    c,
	}
}

//Error return error text
func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %s desc = %s", Code(e.Code), e.Message)
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c Code, format string, a ...interface{}) *Error {
	return New(c, fmt.Sprintf(format, a...))
}

//StatusCode return http status code in error
func (e *Error) StatusCode() int {
	return HTTPStatusFromCode(e.Code)
}

//WithDetails add detail for error
func (e *Error) WithDetails(details ...Detail) *Error {
	e.Details = details
	return e
}
