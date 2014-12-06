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

// Context是tiny的灵魂，集多种功能于一身。
// 它存储了http.Request,http.ResponseWriter这两个基本的对象，这是HTTP的基础；
// 它的Params属性提供了路由参数，Data属性则可以使得中间件之间可以传输数据；
// 它提供在响应返回文本、JSON的方法，以后会支持更多的媒体类型；
// 更重的是，它的Next()方法让中间件变得非常灵活。
type Context struct {
	Req     *http.Request
	Res     http.ResponseWriter
	Params  map[string]string      // 可获取路由中的参数，如"/users/:id"，可用ctx.Params["id"]获取id
	Data    map[string]interface{} // 用于中间件的数据传输
	handles []Handler
	index   int
}

// 执行下一个Handler
func (ctx *Context) Next() {
	ctx.index++
	if ctx.index < len(ctx.handles) {
		ctx.handles[ctx.index](ctx)
	}
}

func (ctx *Context) Text(status int, text string) {
	ctx.Res.Header().Set(contentType, appendCharset(contentText, defaultCharset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write([]byte(text))
}

func (ctx *Context) Json(status int, v interface{}, configs ...map[string]interface{}) error {

	indent := true
	charset := defaultCharset

	if len(configs) == 1 {
		for _, config := range configs {
			if config["indent"].(bool) == false {
				indent = false
			}
			if config["charset"] != nil {
				charset = config["charset"].(string)
			}
		}
	}

	var result []byte
	var err error
	if indent == true {
		result, err = json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
	} else {
		result, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	ctx.Res.Header().Set(contentType, appendCharset(contentJSON, charset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write(result)
	return nil
}

func appendCharset(content string, charset string) string {
	return content + "; charset=" + charset
}
