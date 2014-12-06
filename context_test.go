package tiny

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Context_Text(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := &Context{req, res, make(map[string]string), make(map[string]interface{}), defaultCharset, []Handler{}, 0}
	ctx.Text(200, "hello,world")
	expect(t, res.Code, 200)
	expect(t, res.Header().Get(contentType), appendCharset(contentText, defaultCharset))
	expect(t, res.Body.String(), "hello,world")
}

func Test_Context_Json(t *testing.T) {
	jsonTemplate := `{
  "id": 1,
  "name": "test"
}`
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := &Context{req, res, make(map[string]string), make(map[string]interface{}), defaultCharset, []Handler{}, 0}
	ctx.Json(403, map[string]interface{}{
		"id":   1,
		"name": "test",
	})
	expect(t, res.Code, 403)
	expect(t, res.Header().Get(contentType), appendCharset(contentJSON, defaultCharset))
	expect(t, res.Body.String(), jsonTemplate)
}
