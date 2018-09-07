package ctxrouter

import (
	"encoding/json"
	"github.com/ti/ctxrouter/errors"
	"net/http"
	"reflect"
)

//Params get params form request (It is faster than most other function, because there is no extra compute )
//req http.Request
func Params(req *http.Request) []string {
	return req.Header[paramHeader]
}

const paramHeader = "X-Ctxrouter-Params"

//ContextInterface the interface of any context
//the context must have Init and DecodeRequest
type ContextInterface interface {
	Init(http.ResponseWriter, *http.Request)
	DecodeRequest() error
}

//New new router
func New() *Router {
	return &Router{
		handlers: make(map[string][]Handler),
	}
}

//ServeHTTP just used by system http handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	val, _, params, err := r.Match(req.Method, req.URL.Path)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	if val.callT == nil {
		req.Header[paramHeader] = params
		if h, ok := val.callV.Interface().(http.HandlerFunc); ok {
			h.ServeHTTP(w, req)
		} else if hf, ok := val.callV.Interface().(func(http.ResponseWriter, *http.Request)); ok {
			hf(w, req)
		}
		return
	}
	ctx := reflect.New(val.callT).Interface().(ContextInterface)
	ctx.Init(w, req)
	if err := ctx.DecodeRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	in := []reflect.Value{reflect.ValueOf(ctx)}
	if val.hasParams {
		in = append(in, val.paramsV...)
	}
	rets := val.callV.Call(in)
	var statusError Error
	var data interface{}
	var dataOK bool
	if len(rets) == 1 {
		statusError = errorFromValue(rets[0])
		if statusError == nil && !rets[0].IsNil() {
			data, dataOK = rets[0].Interface().(interface{})
		}
	} else if len(rets) == 2 {
		statusError = errorFromValue(rets[1])
		if statusError == nil && !rets[0].IsNil() {
			data, dataOK = rets[0].Interface().(interface{})
		}
	} else {
		return
	}
	if statusError == nil {
		if dataOK {
			if d, err := json.Marshal(data); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.Write(d)
				return
			}
		}
		return
	}
	d, _ := json.Marshal(statusError)
	w.Header().Set("Content-Type", "application/json")
	statusCode := statusError.StatusCode()
	if statusCode > 0 {
		w.WriteHeader(statusCode)
	} else {
		w.WriteHeader(400)
	}
	w.Write(d)
}

//errorFromValue bool is if the error is nil
func errorFromValue(v reflect.Value) Error {
	if v.IsNil() {
		return nil
	}
	if e, ok := v.Interface().(Error); ok {
		if e.IsNil() {
			return nil
		}
		return e
	}
	if e, ok := v.Interface().(error); ok {
		if e != nil {
			errStr := e.Error()
			if len(errStr) > 0 {
				return errors.CodeError(errors.Unknown).WithDescription(errStr)
			}
		}
	}
	return nil
}

//Get http Get method
func (r *Router) Get(path string, controller interface{}) {
	if err := r.Handle("GET", path, controller); err != nil {
		panic(err)
	}
}

//Post http Post method
func (r *Router) Post(path string, controller interface{}) {
	if err := r.Handle("POST", path, controller); err != nil {
		panic(err)
	}
}

//Patch http Patch method
func (r *Router) Patch(path string, controller interface{}) {
	if err := r.Handle("PATCH", path, controller); err != nil {
		panic(err)
	}
}

//Put http Put method
func (r *Router) Put(path string, controller interface{}) {
	if err := r.Handle("PUT", path, controller); err != nil {
		panic(err)
	}
}

//Delete http Delete method
func (r *Router) Delete(path string, controller interface{}) {
	if err := r.Handle("DELETE", path, controller); err != nil {
		panic(err)
	}
}

//Head http Head method
func (r *Router) Head(path string, controller interface{}) {
	if err := r.Handle("HEAD", path, controller); err != nil {
		panic(err)
	}
}

//Options http Options method
func (r *Router) Options(path string, controller interface{}) {
	if err := r.Handle("OPTIONS", path, controller); err != nil {
		panic(err)
	}
}

//All http all method
func (r *Router) All(path string, controller interface{}) {
	if err := r.Handle("*", path, controller); err != nil {
		panic(err)
	}
}
