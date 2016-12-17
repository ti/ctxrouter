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
	"reflect"
	"strconv"
	"errors"
)


//Param get params form request (It is faster than most other function, because there is no extra compute )
//req http.Request
func Params(req *http.Request) []string {
	return  req.Header[paramHeader]
}

//New new http router, it can be handler by system
func New() *Router {
	return  &Router{tree: &node{}, }
}

const paramHeader  = "X-Ctxrouter-Params"

type (
	//Value your can add any value on it, if you add a interface{} it will go to Value.V
	Value  struct {
		CallV     reflect.Value
		CallT     reflect.Type
		V         interface{}
		ParamsV   []reflect.Value
		ParamsT   []reflect.Type
		HasParams bool //faster when callback
	}
	//ContextInterface you can add anycontext you want if it implement ContextInterface
	ContextInterface interface {
		Init(http.ResponseWriter, *http.Request)
		DecodeRequest() error
	}
	Router struct {
		tree   *node
	}
	leaf struct {
		data map[string]Value
	}
)




func (this *Router) Add(path, method string, v interface{}) error {
	if method == "" {
		method = "default"
	}
	val := Value{
		V:v,
		CallV:reflect.ValueOf(v),
	}
	if reflect.TypeOf(v).Kind() == reflect.Func {
		if _, ok := val.CallV.Interface().(http.HandlerFunc); ok {
			//do noting
		} else if _, ok := val.CallV.Interface().(func(http.ResponseWriter, *http.Request)); ok {
			//do noting
		} else {
			val.CallT = reflect.TypeOf(v).In(0).Elem()
			paramsLen := val.CallV.Type().NumIn()
			val.HasParams = paramsLen > 1
			for i := 0; i < paramsLen; i++ {
				if i > 0 {
					if i == 1 {
						val.ParamsT = make([]reflect.Type, 0)
						val.ParamsT = append(val.ParamsT, val.CallV.Type().In(i))
					} else if i > 1 {
						val.ParamsT = append(val.ParamsT, val.CallV.Type().In(i))
					}
				}
			}
		}
	}
	if this.tree == nil {
		this.tree = &node{}
	}
	if vMap, _, _ := this.tree.getValue(path); vMap != nil {
		if vMap, ok := vMap.(*leaf); ok {
			vMap.data[method] = val
			return nil
		} else {
			panic("router value is not a value map")
		}
	}
	if err := this.tree.addRoute(path, &leaf{data: map[string]Value{method:val}}); err != nil {
		return err
	}
	return nil
}


func (this *Router) Match(method, path string) (val Value, p []string) {
	if v, p, _ := this.tree.getValue(path); v != nil {
		if v, ok := v.(*leaf); ok {
			if v.data[method].V != nil {
				val = v.data[method]
			} else {
				val = v.data["default"]
			}
			if val.V != nil && val.CallT != nil && p != nil {
				val.ParamsV = make([]reflect.Value, 0)
				for i, n := range p {
					pt := val.ParamsT[i]
					pv, err := strConv(n, pt)
					if err == nil {
						val.ParamsV = append(val.ParamsV, pv)
					} else {
						return Value{}, nil
					}
				}
			}
			return val, p
		}
		panic("router value is not valueMap")
	}
	return val, p
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	val, params := this.Match(r.Method, r.URL.Path)
	if val.V == nil {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	if val.CallT == nil {
		r.Header[paramHeader] = params
		if h, ok := val.CallV.Interface().(http.HandlerFunc); ok {
			h.ServeHTTP(w,r)
		} else if hf, ok := val.CallV.Interface().(func(http.ResponseWriter, *http.Request)); ok {
			hf(w,r)
		}
		return
	}
	ctx := reflect.New(val.CallT).Interface().(ContextInterface)
	ctx.Init(w, r)
	if err := ctx.DecodeRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	in := []reflect.Value{reflect.ValueOf(ctx)}
	if val.HasParams {
		in = append(in, val.ParamsV...)
	}
	val.CallV.Call(in)
}

func (this *Router) Get(path string, controller interface{}) {
	if err := this.Add(path, "GET", controller); err != nil {
		panic(err)
	}
}

func (this *Router) Post(path string, controller interface{}) {
	if err := this.Add(path, "POST", controller); err != nil {
		panic(err)
	}
}

func (this *Router) Patch(path string, controller interface{}) {
	if err := this.Add(path, "PATCH", controller); err != nil {
		panic(err)
	}
}

func (this *Router) Put(path string, controller interface{}) {
	if err := this.Add(path, "PUT", controller); err != nil {
		panic(err)
	}
}

func (this *Router) Delete(path string, controller interface{}) {
	if err := this.Add(path, "DELETE", controller); err != nil {
		panic(err)
	}
}

func (this *Router) Head(path string, controller interface{}) {
	if err := this.Add(path, "HEAD", controller); err != nil {
		panic(err)
	}
}

func (this *Router) Options(path string, controller interface{}) {
	if err := this.Add(path, "OPTIONS", controller); err != nil {
		panic(err)
	}
}
func (this *Router) All(path string, controller interface{}) {
	if err := this.Add(path, "", controller); err != nil {
		panic(err)
	}
}
//strConv convert string params to function params
func strConv(src string, t reflect.Type) (rv reflect.Value, err error) {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(src), nil
	case reflect.Int:
		v, err := strconv.Atoi(src)
		if err == nil {
			rv = reflect.ValueOf(v)
		}
		return rv, err
	case reflect.Int64:
		v, err := strconv.ParseInt(src, 10, 64)
		if err == nil {
			rv = reflect.ValueOf(v)
		}
		return rv, err
	case reflect.Bool:
		v, err := strconv.ParseBool(src)
		if err == nil {
			rv = reflect.ValueOf(v)
		}
		return rv, err
	case reflect.Float64:
		v, err := strconv.ParseFloat(src, 64)
		if err == nil {
			rv = reflect.ValueOf(v)
		}
		return rv, err
	case reflect.Float32:
		v, err := strconv.ParseFloat(src, 32)
		if err == nil {
			rv = reflect.ValueOf(float32(v))
		}
		return rv, err
	case reflect.Uint64:
		v, err := strconv.ParseUint(src, 10, 64)
		if err == nil {
			rv = reflect.ValueOf(v)
		}
		return rv, err
	case reflect.Uint32:
		v, err := strconv.ParseUint(src, 10, 32)
		if err == nil {
			rv = reflect.ValueOf(uint32(v))
		}
		return rv, err
	default:
		return rv, errors.New("elem of invalid type")
	}
}