package tiny

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func createReqRes(method string, url string) (req *http.Request, res *httptest.ResponseRecorder) {
	req, _ = http.NewRequest(method, url, nil)
	res = httptest.NewRecorder()
	return
}

func Test_Server(t *testing.T) {
	app := New()

	app.Prepend(func(ctx *Context) {
		ctx.Data["pre1"] = 1
		ctx.Next()
	})

	app.Prepend(func(ctx *Context) {
		ctx.Data["pre2"] = 2
		ctx.Next()
	})

	app.Append(func(ctx *Context) {
		ctx.Data["end1"] = "end1"
		ctx.Next()
	})

	app.Append(func(ctx *Context) {
		ctx.Data["end2"] = "end2"
		ctx.Next()
	})

	app.ErrorHandle(func(ctx *Context) {
		ctx.Text(500, ctx.Data["error"].(string))
	})

	app.NotFound(func(ctx *Context) {
		ctx.Text(404, "not found")
	})

	app.Get("/users/:id", func(ctx *Context) {
		ctx.Json(200, map[string]interface{}{
			"id":   ctx.Params["id"],
			"data": ctx.Data,
		})
	})

	app.Post("/users/:id", func(ctx *Context) {
		ctx.Data["uid"] = ctx.Params["id"]
		ctx.Next()
	}, func(ctx *Context) {
		ctx.Text(201, "created, uid is "+ctx.Data["uid"].(string))
	})

	app.Put("/users/:id/name", func(ctx *Context) {
		ctx.Text(200, "what's your name?")
	})

	app.Delete("/users/:id", func(ctx *Context) {
		ctx.Text(204, "No Content")
	})

	req, res := createReqRes("GET", "/users/123")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
	expect(t, res.Header().Get(contentType), contentJSON+appendCharset)
	expect(t, res.Body.String(), `{
  "data": {
    "pre1": 1,
    "pre2": 2
  },
  "id": "123"
}`)

	req, res = createReqRes("POST", "/users/abc")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 201)
	expect(t, res.Header().Get(contentType), contentText+appendCharset)
	expect(t, res.Body.String(), "created, uid is abc")

	req, res = createReqRes("PUT", "/users/123/name")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
	expect(t, res.Header().Get(contentType), contentText+appendCharset)
	expect(t, res.Body.String(), "what's your name?")

	req, res = createReqRes("DELETE", "/users/123")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 204)
	expect(t, res.Header().Get(contentType), contentText+appendCharset)
	expect(t, res.Body.String(), "No Content")

}
