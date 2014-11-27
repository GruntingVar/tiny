package tiny

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

func Gzip(ctx *Context) {
	var writer http.ResponseWriter
	if !strings.Contains(ctx.Req.Header.Get("Accept-Encoding"), "gzip") {
		ctx.Next()
		return
	} else {
		gz := gzip.NewWriter(ctx.Res)
		defer gz.Close()
		writer = gzipResponseWriter{Writer: gz, ResponseWriter: ctx.Res}
		ctx.Res = writer
		ctx.Res.Header().Set("Content-Encoding", "gzip")
		ctx.Next()
		return
	}
}
