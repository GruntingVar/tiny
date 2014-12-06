package tiny

import (
	"testing"
)

func Test_MiddleWare_Gzip(t *testing.T) {
	var str string
	for i := 0; i < 1000; i++ {
		str += "a"
	}

	app := New()
	app.Use(Gzip)
	app.Get("/", func(ctx *Context) {
		ctx.Text(200, str)
	})

	req, res := createReqRes("GET", "/")
	req.Header.Set("Accept-Encoding", "gzip")
	app.ServeHTTP(res, req)
	testItem(t, res.Header().Get("Content-Encoding"), "gzip", "gzip Content-Encoding")

	appNoGzip := New()
	appNoGzip.Get("/", func(ctx *Context) {
		ctx.Text(200, str)
	})
	req1, res1 := createReqRes("GET", "/")
	appNoGzip.ServeHTTP(res1, req1)

	if len(res.Body.String()) >= len(res1.Body.String()) {
		t.Error("gzip not effect")
	}

	req, res = createReqRes("GET", "/")
	app.ServeHTTP(res, req)
	testNotEqualItem(t, res.Header().Get("Content-Encoding"), "gzip", "gzip Content-Encoding")

}
