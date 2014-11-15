package tiny

import (
	"net/http"
)

type Context struct {
	Req  *http.Request
	Res  http.ResponseWriter
	Data map[string]interface{}
	next bool
}
