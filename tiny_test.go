package tiny

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expect(t *testing.T, actual interface{}, expect interface{}) {
	if actual != expect {
		t.Errorf("expect: %v (type: %v) but got: %v (type: %v)\n", expect, reflect.TypeOf(expect), actual, reflect.TypeOf(actual))
	}
}

func testItem(t *testing.T, actual interface{}, expect interface{}, description string) {
	if actual != expect {
		t.Errorf("%s:\nexpect: %v (type: %v) but got: %v (type: %v)\n", description, expect, reflect.TypeOf(expect), actual, reflect.TypeOf(actual))
	}
}

func testNotEqualItem(t *testing.T, actual interface{}, expect interface{}, description string) {
	if actual == expect {
		t.Errorf("%s:\nexpect: %v (type: %v) should not be: %v (type: %v)\n", description, expect, reflect.TypeOf(expect), actual, reflect.TypeOf(actual))
	}
}

func createReqRes(method string, url string) (req *http.Request, res *httptest.ResponseRecorder) {
	req, _ = http.NewRequest(method, url, nil)
	res = httptest.NewRecorder()
	return
}

func createRestServer(app *Tiny, srcName string) {
	handler := func(ctx *Context) {
		ctx.Res.WriteHeader(200)
	}
	app.Post("/"+srcName, handler)
	app.Get("/"+srcName+"/:id", func(ctx *Context) {
		ctx.Text(200, ctx.Params["id"])
	})
	app.Put("/"+srcName+"/:id", handler)
	app.Patch("/"+srcName+"/:id", handler)
	app.Delete("/"+srcName+"/:id", handler)
	app.Options("/"+srcName+"/:id", handler)
	app.Head("/"+srcName+"/:id", handler)
	app.All("/"+srcName+"/:id", handler)
}

func createComplexServer() *Tiny {
	app := New()
	createRestServer(app, "users")
	createRestServer(app, "blogs")
	createRestServer(app, "images")

	app.Use(func(ctx *Context) {
		ctx.Data["pre1"] = 1
		ctx.Next()
	})

	app.Use(func(ctx *Context) {
		ctx.Data["pre2"] = 2
		ctx.Next()
	})

	app.Use(func(ctx *Context) {
		ctx.Next()
		ctx.Data["after1"] = 3
	})

	app.Get("/", func(ctx *Context) {
		ctx.Res.WriteHeader(200)
	})

	app.All("/all", func(ctx *Context) {
		ctx.Res.WriteHeader(200)
	})

	app.Get("/json", func(ctx *Context) {
		ctx.Json(200, map[string]interface{}{
			"name": "json",
			"data": "jsonData",
		})
	})

	app.Get("/text", func(ctx *Context) {
		ctx.Text(200, "text")
	})

	app.Get("/data", func(ctx *Context) {
		ctx.Json(200, ctx.Data)
	})

	app.Get("/multi", func(ctx *Context) {
		ctx.Next()
	}, func(ctx *Context) {
		ctx.Text(200, "multi")
	})

	app.Get("/panic", func(ctx *Context) {
		panic("panic")
	})

	app.Get("/deep/:did/complex/:cid/path", func(ctx *Context) {
		ctx.Json(200, ctx.Params)
	})

	return app
}

func Test_Tiny_SingleServer(t *testing.T) {
	app := New()
	app.Get("/", func(ctx *Context) {
		ctx.Res.WriteHeader(200)
	})
	req, res := createReqRes("GET", "/")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
}

func Test_Tiny_Run(t *testing.T) {
	app := New()
	go app.Run("3000")
}

func Test_Tiny_Rest(t *testing.T) {
	app := New()
	createRestServer(app, "users")

	req, res := createReqRes("POST", "/users")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("GET", "/users/123")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("PUT", "/users/123")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("PATCH", "/users/123")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("DELETE", "/users/123")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("HEAD", "/users/123")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("OPTIONS", "/users/123")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}
}

func Test_Tiny_AllMethod(t *testing.T) {
	app := createComplexServer()

	req, res := createReqRes("POST", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("GET", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("PUT", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("PATCH", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("DELETE", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("HEAD", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}

	req, res = createReqRes("OPTIONS", "/all")
	app.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fail()
	}
}

func Test_Tiny_MiddleWare(t *testing.T) {
	app := createComplexServer()

	jsonData := `{
  "pre1": 1,
  "pre2": 2
}`
	req, res := createReqRes("GET", "/data")
	app.ServeHTTP(res, req)
	expect(t, res.Body.String(), jsonData)
}

func Test_Tiny_Params(t *testing.T) {
	app := createComplexServer()

	paramData := `{
  "cid": "456",
  "did": "123"
}`
	req, res := createReqRes("GET", "/deep/123/complex/456/path")
	app.ServeHTTP(res, req)
	expect(t, res.Body.String(), paramData)
}

func Test_Tiny_DefaultHandler(t *testing.T) {
	app := createComplexServer()

	req, res := createReqRes("GET", "/not/exists")
	app.ServeHTTP(res, req)
	testItem(t, res.Code, 404, "get not exists path")

	req, res = createReqRes("GET", "/panic")
	app.ServeHTTP(res, req)
	testItem(t, res.Code, 500, "get panic path")
}

func Test_Tiny_CustomHandler(t *testing.T) {
	app := createComplexServer()

	app.NotFound(func(ctx *Context) {
		ctx.Text(200, "custom")
	})

	app.PanicHandler(func(ctx *Context) {
		ctx.Text(200, "don't worry")
	})

	req, res := createReqRes("GET", "/not/exists")
	app.ServeHTTP(res, req)
	testItem(t, res.Code, 200, "get not exists path")
	testItem(t, res.Body.String(), "custom", "get not exists path")

	req, res = createReqRes("GET", "/panic")
	app.ServeHTTP(res, req)
	testItem(t, res.Code, 200, "get panic path")
	testItem(t, res.Body.String(), "don't worry", "get panic path")
}

func Test_Tiny_Handlers(t *testing.T) {
	app := createComplexServer()
	req, res := createReqRes("GET", "/multi")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
}

func Test_Tiny_Text(t *testing.T) {
	app := createComplexServer()
	req, res := createReqRes("GET", "/text")
	app.ServeHTTP(res, req)
	expect(t, res.Code, 200)
	expect(t, res.Body.String(), "text")
}

func Test_Tiny_Json(t *testing.T) {
	app := createComplexServer()
	req, res := createReqRes("GET", "/json")
	app.ServeHTTP(res, req)

	jsonData := `{
  "data": "jsonData",
  "name": "json"
}`
	expect(t, res.Code, 200)
	expect(t, res.Body.String(), jsonData)
}
