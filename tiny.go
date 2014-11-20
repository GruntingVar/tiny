package tiny

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

var defaultNotFoundHandle = func(ctx *Context) {
	ctx.Text(404, "Not Found")
}

var defaultErrorHandle = func(ctx *Context) {
	ctx.Text(500, "Internal Server Error")
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

type Tiny struct {
	root           *routeNode
	preHandles     []Handle // 在进行路由匹配之前执行的handle
	endHandles     []Handle // 在进行路由匹配之后执行的handle
	notFoundHandle Handle
	errorHandle    Handle
}

func New() *Tiny {
	return &Tiny{
		root:           createRoot(),
		preHandles:     []Handle{},
		endHandles:     []Handle{},
		notFoundHandle: defaultNotFoundHandle,
		errorHandle:    defaultErrorHandle,
	}
}

func (tiny *Tiny) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	found, node, data := tiny.root.findUrl(r.URL.Path)

	var writer http.ResponseWriter
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		writer = rw
	} else {
		rw.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(rw)
		defer gz.Close()
		writer = gzipResponseWriter{Writer: gz, ResponseWriter: rw}
	}
	ctx := &Context{r, writer, data, make(map[string]interface{}), false}

	defer func(ctx *Context) {
		if err := recover(); err != nil {
			ctx.Data["error"] = err
			tiny.errorHandle(ctx)
		}
	}(ctx)

	runHandles(ctx, tiny.preHandles)
	if found == true {
		handles := node.getHandles(strings.ToUpper(r.Method))
		if handles == nil {
			tiny.notFoundHandle(ctx)
		} else {
			runHandles(ctx, handles)
		}
	} else {
		tiny.notFoundHandle(ctx)
	}

	runHandles(ctx, tiny.endHandles)
}

func (tiny *Tiny) Run(port string) {
	log.Println("Tiny is listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, tiny))
}

func (tiny *Tiny) Prepend(h Handle) {
	tiny.preHandles = append(tiny.preHandles, h)
}

func (tiny *Tiny) Append(h Handle) {
	tiny.endHandles = append(tiny.endHandles, h)
}

func (tiny *Tiny) NotFound(h Handle) {
	tiny.notFoundHandle = h
}

func (tiny *Tiny) ErrorHandle(h Handle) {
	tiny.errorHandle = h
}

func (tiny *Tiny) Post(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Post(handles)
}

func (tiny *Tiny) Get(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Get(handles)
}

func (tiny *Tiny) Put(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Put(handles)
}

func (tiny *Tiny) Patch(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Patch(handles)
}

func (tiny *Tiny) Delete(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Delete(handles)
}

func (tiny *Tiny) Head(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Head(handles)
}

func (tiny *Tiny) Options(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.Options(handles)
}

func (tiny *Tiny) All(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.All(handles)
}

func runHandles(ctx *Context, handles []Handle) {
	ctx.next = false
	for _, handle := range handles {
		handle(ctx)
		if ctx.next == false {
			return
		}
	}
}
