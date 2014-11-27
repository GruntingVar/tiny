package tiny

import (
	"testing"
)

func Test_Gzip(t *testing.T) {
	app := New()

	app.Use(Gzip)

	app.Get("/", func(ctx *Context) {
		ctx.Text(200, "gzip test")
	})

	req, res := createReqRes("GET", "/")
	req.Header.Set("Accept-Encoding", "gzip")
	app.ServeHTTP(res, req)
	expect(t, res.Header().Get("Content-Encoding"), "gzip")
}
