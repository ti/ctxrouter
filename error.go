package ctxrouter

import (
	"encoding/json"
	"net/http"
	"strings"
)

//Error Default HTTP Error
type Error struct {
	Status           int         `json:"-"`
	ErrorType        string      `json:"error,omitempty"`
	ErrorDescription string      `json:"error_description,omitempty"`
	ErrorCode        int         `json:"error_code,omitempty"`
	ErrorUri         string      `json:"error_uri,omitempty"`
	ErrorData        interface{} `json:"error_data,omitempty"`
}

//Error return error text
func (this *Error) Error() string {
	resp := this.ErrorType
	if this.ErrorDescription != "" {
		resp += ": " + this.ErrorDescription
	}
	return resp
}

//StatusCode return http status code in error
func (this *Error) StatusCode() int {
	return this.Status
}

//SetDescription the error detail
func (this *Error) SetDescription(description string) *Error {
	this.ErrorDescription = description
	return this
}

//SetErrorCode set error code
func (this *Error) SetErrorCode(code int) *Error {
	this.ErrorCode = code
	return this
}

//SetErrorDescription set errorDescription by error
func (this *Error) SetErrorDescription(errorDescription error) *Error {
	if errorDescription != nil {
		this.ErrorDescription = errorDescription.Error()
	}
	return this
}

//SetUri set error uri
func (this *Error) SetUri(uri string) *Error {
	this.ErrorUri = uri
	return this
}

//SetStatus set http status code in error
func (this *Error) SetStatus(status int) *Error {
	this.Status = status
	return this
}

//SetData set error data
func (this *Error) SetData(data interface{}) *Error {
	this.ErrorData = data
	return this
}

//NewError new error by type
func NewError(t string) *Error {
	return &Error{Status: 400, ErrorType: t}
}

//HttpStatusError new error by http status code
//the error type is http status text
func HttpStatusError(status int) *Error {
	return &Error{Status: status, ErrorType: strings.ToLower(strings.Replace(http.StatusText(status), " ", "_", -1))}
}

//JSONResponse response json to any http writer
//if data is a error it will response a errror json
func JSONResponse(w http.ResponseWriter, data interface{}) {
	if err, ok := data.(*Error); ok {
		JSONResponseVerbose(w, err.Status, nil, err)
	} else {
		JSONResponseVerbose(w, 200, nil, data)
	}
}

//JSONResponseVerbose response json to any http writer, include status ,header, data
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
