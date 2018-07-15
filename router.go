package ctxrouter

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

//Router the router
type Router struct {
	handlers map[string][]Handler
}

//Handle handler path in router
func (s *Router) Handle(method, path string, v interface{}) error {
	path = adapterRouterStyle(path)
	pattern, err := ParsePatternURL(path)
	if err != nil {
		return err
	}
	val := Handler{
		V:     v,
		Pat:   pattern,
		callV: reflect.ValueOf(v),
	}
	if reflect.TypeOf(v).Kind() == reflect.Func {
		switch val.callV.Interface().(type) {
		case http.HandlerFunc, func(http.ResponseWriter, *http.Request):
			//do noting
		default:
			val.callT = reflect.TypeOf(v).In(0).Elem()
			paramsLen := val.callV.Type().NumIn()
			val.hasParams = paramsLen > 1
			for i := 0; i < paramsLen; i++ {
				if i == 1 {
					val.paramsT = make([]reflect.Type, 0)
					val.paramsT = append(val.paramsT, val.callV.Type().In(i))
				} else if i > 1 {
					val.paramsT = append(val.paramsT, val.callV.Type().In(i))
				}
			}
		}
	}
	if method == "" {
		method = "*"
	}
	s.handlers[method] = append(s.handlers[method], val)
	return nil
}

// Match dispatches the request to the first handler whose pattern matches to r.Method and r.Path.
func (s *Router) Match(method string, path string) (h Handler, pathParams map[string]string, paramsList []string, err error) {
	components := strings.Split(path[1:], "/")
	l := len(components)
	var verb string
	if idx := strings.LastIndex(components[l-1], ":"); idx > 0 {
		c := components[l-1]
		components[l-1], verb = c[:idx], c[idx+1:]
	}
	h, pathParams, paramsList, err = match(s.handlers[method], components, verb)
	if err != nil {
		return match(s.handlers["*"], components, verb)
	}
	return
}

func match(handlers []Handler, components []string, verb string) (h Handler, pathParams map[string]string, paramsList []string, err error) {
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
					return handler, pathParams, p, ErrNotMatch
				}
			}
		}
		return handler, pathParams, p, nil
	}
	err = ErrNotMatch
	return
}

//Handler the handler instance in router
type Handler struct {
	Pat Pattern
	V   interface{}

	//some values for reflect call
	callV   reflect.Value
	callT   reflect.Type
	paramsV []reflect.Value
	paramsT []reflect.Type
	//faster when callback
	hasParams bool
}

//adapterRouterStyle change /v1/home/:id/name style to /v1/home/{id}/name style
//
// Deprecated: use /v1/home/{id}/name style
func adapterRouterStyle(src string) string {
	var prefix bool
	for i := 0; i < len(src); i++ {
		v := src[i]
		if prefix && v == '/' {
			src = src[0:i] + "}" + src[i:]
			prefix = false
			continue
		}
		if !prefix {
			if v == ':' && src[i-1] == '/' {
				src = src[0:i] + "{" + src[i+1:]
				prefix = true
			} else if v == '*' && src[i-1] == '/' && i < len(src)-1 {
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

//ParsePatternURL parse any path to google pattern
func ParsePatternURL(path string) (pattern Pattern, err error) {
	cp, err := Parse(path)
	if err != nil {
		return
	}
	tp := cp.Compile()
	pattern, err = NewPattern(tp.OpCodes, tp.Pool, tp.Verb)
	return
}

//strConv convert string params to function params
func strConv(src string, t reflect.Type) (rv reflect.Value, err error) {
	switch t.Kind() {
	case reflect.String:
		rv = reflect.ValueOf(src)
	case reflect.Int:
		v, err := strconv.Atoi(src)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(v)
		})
	case reflect.Int64:
		v, err := strconv.ParseInt(src, 10, 64)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(v)
		})
	case reflect.Bool:
		v, err := strconv.ParseBool(src)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(v)
		})
	case reflect.Float64:
		v, err := strconv.ParseFloat(src, 64)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(v)
		})
	case reflect.Float32:
		v, err := strconv.ParseFloat(src, 32)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(float32(v))
		})
	case reflect.Uint64:
		v, err := strconv.ParseUint(src, 10, 64)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(v)
		})
	case reflect.Uint32:
		v, err := strconv.ParseUint(src, 10, 32)
		noErrorFunc(err, func() {
			rv = reflect.ValueOf(uint32(v))
		})
	default:
		err = errors.New("elem of invalid type")
	}
	return rv, err
}

func noErrorFunc(err error, fn func()) {
	if err != nil {
		return
	}
	fn()
}
