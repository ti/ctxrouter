package ctxrouter

import (
	"strings"
	"reflect"
	"net/http"
	"strconv"
	"errors"
)

type Router struct {
	handlers               map[string][]Handler
}

func (s *Router) Handle(method, path string, v interface{}) error {
	path = adapterRouterStyle(path)
	pattern, err := ParsePatternUrl(path)
	if err != nil {
		return err
	}
	val := Handler{
		V:v,
		Pat:pattern,
		callV:reflect.ValueOf(v),
	}
	if reflect.TypeOf(v).Kind() == reflect.Func {
		if _, ok := val.callV.Interface().(http.HandlerFunc); ok {
			//do noting
		} else if _, ok := val.callV.Interface().(func(http.ResponseWriter, *http.Request)); ok {
			//do noting
		} else {
			val.callT = reflect.TypeOf(v).In(0).Elem()
			paramsLen := val.callV.Type().NumIn()
			val.hasParams = paramsLen > 1
			for i := 0; i < paramsLen; i++ {
				if i > 0 {
					if i == 1 {
						val.paramsT = make([]reflect.Type, 0)
						val.paramsT = append(val.paramsT, val.callV.Type().In(i))
					} else if i > 1 {
						val.paramsT = append(val.paramsT, val.callV.Type().In(i))
					}
				}
			}
		}
	}
	s.handlers[method] = append(s.handlers[method], val)
	return nil
}


// Match dispatches the request to the first handler whose pattern matches to r.Method and r.Path.
func (s *Router) Match(method string, path string) (h Handler, pathParams map[string]string, paramsList []string, err error){
	components := strings.Split(path[1:], "/")
	l := len(components)
	var verb string
	if idx := strings.LastIndex(components[l-1], ":"); idx > 0 {
		c := components[l-1]
		components[l-1], verb = c[:idx], c[idx+1:]
	}
	handlers, ok := s.handlers[method]
	if !ok {
		handlers, ok = s.handlers["*"]
	}
	if ok {
		for _, handler := range handlers {
			pathParams, p, err := handler.Pat.Match(components, verb)
			if err != nil {
				continue
			}
			if handler.V != nil && handler.callT != nil && p != nil && len(p) == len(handler.paramsT) {
				handler.paramsV = make([]reflect.Value, 0)
				for i, n := range p {
					pt := handler.paramsT[i]
					pv, err := strConv(n, pt)
					if err == nil {
						handler.paramsV = append(handler.paramsV, pv)
					} else {
						return handler, pathParams, p,ErrNotMatch
					}
				}
			}
			return handler, pathParams, p,nil
		}
	}
	err = ErrNotMatch
	return
}

type Handler struct {
	Pat Pattern
	V   interface{}

	//some values for reflect call
	callV     reflect.Value
	callT     reflect.Type
	paramsV   []reflect.Value
	paramsT   []reflect.Type
	//faster when callback
	hasParams bool
}



//adapterStr change /v1/home/:id/name style to /v1/home/{id}/name style
//
// Deprecated: use /v1/home/{id}/name style
func adapterRouterStyle(src string) string {
	var prefix bool
	for i :=0; i< len(src); i ++ {
		v := src[i]
		if prefix && v == '/' {
			src = src[0:i] + "}" +  src[i:]
			prefix = false
			continue
		}
		if !prefix {
			if v == ':' && src[i-1] == '/' {
				src = src[0:i] + "{" + src[i+1:]
				prefix = true
			} else if v == '*' && src[i-1] == '/' && i < len(src) -1 {
				src = src[0:i] + "{" + src[i+1:] + "=**}"
				break
			}
		}
	}
	if prefix {
		src += "}"
	}
	return src
}


func ParsePatternUrl(path string)  (pattern Pattern, err error){
	cp, err := Parse(path)
	if err != nil {
		return
	}


	tp := cp.Compile()


	if strings.HasSuffix(path,"/_") {
		 tp.Pool[len(tp.Pool) - 1 ] = ""
	}
	pattern , err = NewPattern(tp.OpCodes, tp.Pool, tp.Verb)
	return
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