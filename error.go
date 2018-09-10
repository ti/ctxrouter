package ctxrouter

import (
	"encoding/json"
	"net/http"
)

//Error You can custom any error structure you want
type Error interface {
	StatusCode() int
	Error() string
	IsNil() bool
}

//JSONResponse response json to any http writer
//if data is a error it will response a errror json
func JSONResponse(w http.ResponseWriter, data interface{}) {
	if err, ok := data.(Error); ok {
		JSONResponseVerbose(w, err.StatusCode(), nil, err)
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
