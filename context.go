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
	Params  map[string]string
	Data    map[string]interface{}
	handles []Handler
	index   int
}

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
