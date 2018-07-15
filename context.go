package ctxrouter

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

//Context the context of http
//you can use by startContext, nextContext, procContext, endContext ....
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Data    interface{}
}

//Init the start of context
//you can define you own link (c *context){c.Context.Init(w,r), your code ...}
func (c *Context) Init(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
}

//DecodeRequest You can implement your DecodeRequest, it can be form or something else
func (c *Context) DecodeRequest() error {
	if c.Data != nil && strings.Contains(c.Request.Header.Get("Content-Type"), "json") {
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&c.Data); err != nil {
			return errors.New("json decode error - " + err.Error())
		}
		return nil
	}
	return nil
}

//DecodeJson decode json
func (c *Context) DecodeJson(data interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(data); err != nil {
		return errors.New("json decode error - " + err.Error())
	}
	return nil
}

//JSON response json
func (c *Context) JSON(data interface{}) {
	if d, err := json.Marshal(data); err != nil {
		panic("Error marshalling json: %v:" + err.Error())
	} else {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(d)
	}
}

//Text response textplain
func (c *Context) Text(data string) {
	io.WriteString(c.Writer, data)
}

//Redirect http 302 to url
func (c *Context) Redirect(urlStr string, code int) {
	http.Redirect(c.Writer, c.Request, urlStr, code)
}

//Status set response status code
func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
}

//StatusText response textplain by http.Status code
func (c *Context) StatusText(status int) {
	io.WriteString(c.Writer, http.StatusText(status))
}

//StatusError output standard error json body by http.Status code
//exp: StatusError(404,"not fond something"),will response {"error":"not_found", "error_description":"not fond something"}
func (c *Context) StatusError(status int, errorDescription string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(`{"error":"` + strings.ToLower(strings.Replace(http.StatusText(status), " ", "_", -1)) + `","error_description":"` + errorDescription + `"}`))
}
