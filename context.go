package tiny

import (
	"encoding/json"
	"net/http"
)

const (
	contentType    = "Content-Type"
	contentText    = "text/plain"
	contentJSON    = "application/json"
	defaultCharset = "UTF-8"
	appendCharset  = "; charset=" + defaultCharset
)

type Context struct {
	Req    *http.Request
	Res    http.ResponseWriter
	Params map[string]string
	Data   map[string]interface{}
	next   bool
}

func (ctx *Context) Next() {
	ctx.next = true
}

func (ctx *Context) Text(status int, text string) {
	ctx.Res.Header().Set(contentType, contentText+appendCharset)
	ctx.Res.WriteHeader(status)
	ctx.Res.Write([]byte(text))
}

func (ctx *Context) Json(status int, v interface{}) error {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	ctx.Res.Header().Set(contentType, contentJSON+appendCharset)
	ctx.Res.WriteHeader(status)
	ctx.Res.Write(result)
	return nil
}
