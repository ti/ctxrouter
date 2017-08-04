package ctxrouter

import (
	"net/http"
	"strings"
	"encoding/json"
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

func (this *Error) SetData(data interface{}) *Error {
	this.ErrorData = data
	return this
}


func NewError(t string) *Error {
	return &Error{Status: 400, ErrorType: t}
}


func HttpStatusError(status int) *Error {
	return &Error{Status: status, ErrorType: strings.ToLower(strings.Replace(http.StatusText(status), " ", "_", -1))}
}


func JSONResponse(w http.ResponseWriter, data interface{}) {
	if err, ok := data.(*Error); ok {
		JSONResponseVerbose(w, err.Status, nil, err)
	} else {
		JSONResponseVerbose(w, 200, nil, data)
	}
}

func JSONResponseVerbose(w http.ResponseWriter, status int, header http.Header, data interface{}) {
	if header != nil {
		for k, v := range header {
			for _, vv := range v {
				w.Header().Set(k, vv)
			}
		}
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.Header().Del("Content-Length")

	if bs, ok := data.([]byte); ok {
		w.WriteHeader(status)
		w.Write(bs)
		return
	}
	if d, err := json.Marshal(data); err != nil {
		panic("Error marshalling json: %v:" + err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		w.Write(d)
		return
	}
}