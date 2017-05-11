package ctxrouter

import (
	"net/http"
	"strings"
)

type Error struct {
	Status           int    `json:"-"`
	ErrorType        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
	ErrorData        interface{} `json:"error_data,omitempty"`
}

func (this *Error) Error() string {
	resp := this.ErrorType
	if (this.ErrorDescription != "") {
		resp += ": " + this.ErrorDescription
	}
	return resp
}

func (this *Error) SetDescription(description string) *Error {
	this.ErrorDescription = description
	return this
}

func (this *Error) SetErrorDescription(errorDescription error) *Error {
	if errorDescription != nil {
		this.ErrorDescription = errorDescription.Error()
	}
	return this
}

func (this *Error) SetUri(uri string) *Error {
	this.ErrorUri = uri
	return this
}

func (this *Error) SetStatus(status int) *Error {
	this.Status = status
	return this
}

func NewError(t string) *Error {
	return &Error{Status: 400, ErrorType: t}
}



func HttpStatusError(status int) *Error {
	return &Error{Status: status, ErrorType: strings.ToLower(strings.Replace(http.StatusText(status), " ", "_", -1))}
}
