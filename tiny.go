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
	middlewares    []Handle // 在进行路由匹配之前执行的handle
	notFoundHandle Handle
	errorHandle    Handle
}

func New() *Tiny {
	return &Tiny{
		root:           createRoot(),
		middlewares:    []Handle{},
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

	var handles []Handle
	if found == true {
		handles = append(tiny.middlewares, node.getHandles(strings.ToUpper(r.Method))...)
	} else {
		handles = append(tiny.middlewares, tiny.notFoundHandle)
	}

	ctx := &Context{r, writer, data, make(map[string]interface{}), handles, 0}

	defer func(ctx *Context) {
		if err := recover(); err != nil {
			ctx.Data["error"] = err
			tiny.errorHandle(ctx)
		}
	}(ctx)

	handles[0](ctx)
}

func (tiny *Tiny) Run(port string) {
	log.Println("Tiny is listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, tiny))
}

func (tiny *Tiny) Use(h Handle) {
	tiny.middlewares = append(tiny.middlewares, h)
}

func (tiny *Tiny) NotFound(h Handle) {
	tiny.notFoundHandle = h
}

func (tiny *Tiny) ErrorHandle(h Handle) {
	tiny.errorHandle = h
}

func (tiny *Tiny) Post(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.post(handles)
}

func (tiny *Tiny) Get(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.get(handles)
}

func (tiny *Tiny) Put(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.put(handles)
}

func (tiny *Tiny) Patch(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.patch(handles)
}

func (tiny *Tiny) Delete(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.delete(handles)
}

func (tiny *Tiny) Head(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.head(handles)
}

func (tiny *Tiny) Options(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.options(handles)
}

func (tiny *Tiny) All(url string, handles ...Handle) {
	node := tiny.root.addUrl(url)
	node.all(handles)
}
