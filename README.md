# A High performance HTTP request router with Context


# Features

* no regexp （faster）
* wildcards router support
* can decode request before business layer (JSON, xml or other)
* handler simple and pro

# Hello World Sample

```go
package main

import (
	"net/http"
	"fmt"
	"github.com/leenanxi/ctxrouter"
)

func main() {
	r := ctxrouter.New()
	r.Get("/", (*Controller).Index)
	r.Get("/api/:name", (*Controller).Hello)
	//auto decode string and int
	r.Get("/api/:name/json/:age", (*Controller).Json)
	r.All("/api/:name/error", (*Controller).Error)
	r.Get("/api/:name/simple",Simple)
	r.All("/all/*path",All)
	r.Post("/hello",(*UserContext).PrintHello)
	http.ListenAndServe(":8081", r)

}

type Controller struct {
	ctxrouter.Context
}


func (c *Controller) Index() {
	c.Text("index")
}

func (c *Controller) Hello(name string) {
	fmt.Fprintln(c.Writer, "hello "+name)
}

func (c *Controller) Error(name string) {
	c.StatusError(400, name + " is error")
}

func (c *Controller) Json(name string, age int) {
	type Person struct {
		Name string
		Age   int
	}
	c.JSON(Person{Name:name,Age:age})
}

func Simple(ctx *ctxrouter.Context) {
	ctx.Text("simple")
}

func All(ctx *ctxrouter.Context,  path string) {
	ctx.Text("all router goes here " +  path)
}

//decode request
type User struct {
	Id      int              `json:"int"`
	Name    string          `json:"name"`
}


type UserContext struct {
	ctxrouter.Context
	Data  *User
}

//Auto Decode Json or other request if  c.Request.Header.Get("Content-Type") contains json
func (ctx *UserContext) DecodeRequest() error{
	ctx.Data = new(User)
	ctx.Context.Data = ctx.Data
	return ctx.Context.DecodeRequest()
}


func (ctx *UserContext) PrintHello() {
	ctx.Text("Hello "+ ctx.Data.Name)
}

```

# Hello Word Simple Request

```bash
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{"name":"leenanxi"}' \
 'http://localhost:8081/hello'
```


# Thanks 

* tree.go & tree_test.go is edited from httprouter https://github.com/julienschmidt/httprouter