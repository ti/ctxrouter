package main

//server side full context example

import (
	"github.com/ti/ctxrouter"
	"github.com/gorilla/handlers"
	"os"
	"net/http"
	"log"
	"errors"
)

//decode request sample
type User struct {
	Id      int             `json:"id"`
	Name    string          `json:"name"`
}

type UserContext struct {
	ctxrouter.Context
	Data  *User
}

//Auto Decode Json or other request
func (ctx *UserContext) DecodeRequest() error {
	ctx.Data = new(User)
	ctx.Context.Data = ctx.Data
	return ctx.Context.DecodeRequest()
}

//context style
func (ctx *UserContext) Hello(name string) {
	ctx.Text("hello " + name)
}

//resp json data and error
func (ctx *UserContext) RespData(msg string) (interface{}, error) {
	if msg == "error" {
		return nil,ctxrouter.HttpStatusError(400).SetDescription("hello error")
	}
	if msg == "error1" {
		return nil, errors.New("error")
	}
	if msg == "data" {
		return ctx.Data, nil
	}
	return map[string]bool{"success":true},nil
}

func main() {
	r := ctxrouter.New()
	r.Get("/hello/:name", (*UserContext).Hello)
	r.Post("/resp/:msg", (*UserContext).RespData)
	log.Println("server at 8081")
	http.ListenAndServe(":8081", (handlers.LoggingHandler(os.Stdout, r)))
}
