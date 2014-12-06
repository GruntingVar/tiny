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
	ctx.Res.Header().Set(contentType, appendCharset(contentText, defaultCharset))
	ctx.Res.WriteHeader(status)
	ctx.Res.Write([]byte(text))
}

// 向响应写入JSON对象，并设置状态码
//
// status 响应状态码
//
// v 需要转换为Json对象的数据
//
// configs支持如下设置：
// 	map[string]interface{}{
// 		"indent": false, // 是否缩进，如果设为false，则关闭缩进
// 		"charset": "gbk", // 设置编码，默认为UTF-8
// 	}
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
