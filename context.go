package tiny

import (
	"net/http"
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
