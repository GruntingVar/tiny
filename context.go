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
	Req     *http.Request
	Res     http.ResponseWriter
	Params  map[string]string      // 可获取路由中的参数，如"/users/:id"，可用ctx.Params["id"]获取id
	Data    map[string]interface{} // 用于中间件的数据传输
	Charset string                 // Content-Type中的charset，在创建Context实例时会默认设为"UTF-8"，可在调用Text()、Json()之类的方法前更改
	handles []Handler
	index   int
}

// 立即执行下一个handlers
func (ctx *Context) Next() {
	ctx.index++
	if ctx.index < len(ctx.handles) {
		ctx.handles[ctx.index](ctx)
	}
}

// 向响应写入文本，并设置状态码
func (ctx *Context) Text(status int, text string) {
	ctx.Res.Header().Set(contentType, appendCharset(contentText, ctx.Charset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write([]byte(text))
}

// 向响应写入JSON对象，并设置状态码
func (ctx *Context) Json(status int, v interface{}) error {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	ctx.Res.Header().Set(contentType, appendCharset(contentJSON, ctx.Charset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write(result)
	return nil
}

func appendCharset(content string, charset string) string {
	return content + "; charset=" + charset
}
