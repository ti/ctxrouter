package ctxrouter

import (
	"net/http"
	"encoding/json"
	"io"
	"strings"
	"errors"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Data    interface{}
}

func (c *Context) Init(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
}

//you can implement your DecodeRequest, it can be form or something else
func (c *Context) DecodeRequest() error {
	if strings.Contains(c.Request.Header.Get("Content-Type"), "json") {
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&c.Data); err != nil {
			return errors.New("json decode error - " + err.Error())
		}
		return nil
	}
	return nil
}

func (c *Context) JSON(data interface{}) {
	if d, err := json.Marshal(data); err != nil {
		panic("Error marshalling json: %v:" + err.Error())
	} else {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(d)
	}
}

func (c *Context) Text(data string) {
	io.WriteString(c.Writer, data)
}


func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
}

func (c *Context) StatusText(status int) {
	io.WriteString(c.Writer, http.StatusText(status))
}

//exp := write {"error":"not_found", "error_description":"not fond something"}
func (c *Context) StatusError(status int, errorDescription string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(`{"error":"` +  strings.ToLower(strings.Replace(http.StatusText(status), " ", "_", -1)) + `","error_description":"` +  errorDescription + `"}`))
}