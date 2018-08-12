package errors

import (
	"fmt"

	"github.com/ti/ctxrouter/codes"
	detail "github.com/ti/ctxrouter/errordetails"
)

//errorBody Default HTTP Error
type errorBody struct {
	Message string               `json:"error,omitempty"`
	Code    codes.Code           `json:"code,omitempty"`
	Details []detail.ErrorDetail `json:"details,omitempty"`
}

// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *errorBody {
	if c == codes.OK {
		return nil
	}
	return &errorBody{
		Message: msg,
		Code:    c,
	}
}

//Error return error text
func (e *errorBody) Error() string {
	return fmt.Sprintf("error: code = %s desc = %s", codes.Code(e.Code), e.Message)
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c codes.Code, format string, a ...interface{}) *errorBody {
	return New(c, fmt.Sprintf(format, a...))
}

//StatusCode return http status code in error
func (e *errorBody) StatusCode() int {
	return codes.HTTPStatusFromCode(e.Code)
}

//StatusCode return http status code in error
func (e *errorBody) WithDetails(details ...detail.ErrorDetail) *errorBody {
	e.Details = details
	return e
}
