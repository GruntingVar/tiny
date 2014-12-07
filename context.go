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
)

type Context struct {
	Req             *http.Request
	Res             http.ResponseWriter
	Params          map[string]string      // 可获取路由中的参数，如"/users/:id"，可用ctx.Params["id"]获取id
	Data            map[string]interface{} // 用于中间件的数据传输
	Charset         string                 // Content-Type中的charset，在创建Context实例时会默认设为"UTF-8"，可在调用Text()、Json()之类的方法前更改
	PanicMsg        error
	ErrorMsg        error
	middlewares     []Handler // 中间件
	handlers        []Handler // 路由处理的handler
	notfoundHandler Handler   // 没有匹配到路由的handler
	errorHandler    Handler   // 调用Error方法或是使用Text、Json等方法发生错误时会进入该Handler
	middlewareIndex int
	handlersIndex   int
	isMatch         bool
	currentHandler  string // "middleware", "route" or "notfound", 目前在使用的handlers
}

// 立即执行下一个handlers
func (ctx *Context) Next() {
	switch ctx.currentHandler {
	case "middleware":
		ctx.middlewareIndex++
		if ctx.middlewareIndex < len(ctx.middlewares) {
			ctx.middlewares[ctx.middlewareIndex](ctx)
		} else {
			if ctx.isMatch == true {
				ctx.currentHandler = "route"
				ctx.handlers[ctx.handlersIndex](ctx)
			} else {
				ctx.currentHandler = "notfound"
				ctx.notfoundHandler(ctx)
			}
		}
	case "route":
		ctx.handlersIndex++
		if ctx.handlersIndex < len(ctx.handlers) {
			ctx.handlers[ctx.handlersIndex](ctx)
		}
	case "notfound":
		ctx.notfoundHandler(ctx)
	}
}

func (ctx *Context) Error(err error) {
	ctx.ErrorMsg = err
	ctx.errorHandler(ctx)
}

// 向响应写入文本，并设置状态码
func (ctx *Context) Text(status int, text string) {
	ctx.Res.Header().Set(contentType, appendCharset(contentText, ctx.Charset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write([]byte(text))
}

// 向响应写入JSON对象，并设置状态码
func (ctx *Context) Json(status int, v interface{}) {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Res.Header().Set(contentType, appendCharset(contentJSON, ctx.Charset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write(result)
	return
}

func appendCharset(content string, charset string) string {
	return content + "; charset=" + charset
}
