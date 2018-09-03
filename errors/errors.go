package errors

import (
	"encoding/json"
	"fmt"
)

//MarshalJSONWithCode show code field when  MarshalJSON
var MarshalJSONWithCode = false

//Error Default HTTP Error
type Error struct {
	Message string   `json:"error,omitempty"`
	Code    Code     `json:"code,omitempty"`
	Details []Detail `json:"details,omitempty"`
	//compact for simple error
	Description string `json:"error_description,omitempty"`
}

//alias of Error without code json output
type alias struct {
	Message     string   `json:"error,omitempty"`
	Code        Code     `json:"-"`
	Details     []Detail `json:"details,omitempty"`
	Description string   `json:"error_description,omitempty"`
}

//MarshalJSON custom json output
func (e *Error) MarshalJSON() ([]byte, error) {
	if MarshalJSONWithCode {
		return json.Marshal(e)
	}
	return json.Marshal(alias(*e))
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

//CodeError new error by http status code
func CodeError(c Code) *Error {
	return &Error{Code: c, Message: c.String()}
}

//Error return error text
func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("error: code = %d", e.Code)
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

//WithDescription add description for error
func (e *Error) WithDescription(description string) *Error {
	e.Description = description
	return e
}
