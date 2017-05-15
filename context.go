// Copyright 2016 leenanxi All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

//DecodeRequest You can implement your DecodeRequest, it can be form or something else
func (c *Context) DecodeRequest() error {
	if c.Data != nil && strings.Contains(c.Request.Header.Get("Content-Type"), "json"){
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&c.Data); err != nil {
			return errors.New("json decode error - " + err.Error())
		}
		return nil
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

func (c *Context) Redirect(urlStr string, code int) {
	http.Redirect(c.Writer,c.Request,urlStr,code)
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
	c.Writer.Write([]byte(`{"error":"` +  strings.ToLower(strings.Replace(http.StatusText(status), " ", "_", -1)) + `","error_description":"` +  errorDescription + `"}`))
}
