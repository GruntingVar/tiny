// tiny是一个高效、易用、轻量级的用于构建RESTful API Server的框架。
//
// 使用示例：
//
//	package main
//
//	import "github.com/GruntingVar/tiny"
//
//	func main() {
//		app := tiny.New()
//		app.Get("/", func(ctx *tiny.Context) {
//			ctx.Text(200, "Hello, world!")
//		})
//		app.Run("3000")
//	}
//
// 详细文档请访问 https://github.com/GruntingVar/tiny
package tiny

import (
	"log"
	"net/http"
	"strings"
)

var defaultNotFoundHandler = func(ctx *Context) {
	ctx.Text(404, "Not Found")
}

var defaultPanicHandler = func(ctx *Context) {
	ctx.Text(500, "Internal Server Error")
}

// Tiny实现了http.Handler接口，提供HTTP服务
type Tiny struct {
	root            *routeNode
	middlewares     []Handler
	notFoundHandler Handler
	panicHandler    Handler
}

// 创建并返回一个Tiny实例
func New() *Tiny {
	return &Tiny{
		root:            createRoot(),
		middlewares:     []Handler{},
		notFoundHandler: defaultNotFoundHandler,
		panicHandler:    defaultPanicHandler,
	}
}

func (tiny *Tiny) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	found, node, data := tiny.root.findUrl(r.URL.Path)

	var handlers []Handler
	if found == true {
		handlers = append(tiny.middlewares, node.getHandles(strings.ToUpper(r.Method))...)
	} else {
		handlers = append(tiny.middlewares, tiny.notFoundHandler)
	}

	ctx := &Context{r, rw, data, make(map[string]interface{}), handlers, 0}

	defer func(ctx *Context) {
		if err := recover(); err != nil {
			ctx.Data["error"] = err
			tiny.panicHandler(ctx)
		}
	}(ctx)

	handlers[0](ctx)
}

// 监听端口，提供HTTP服务
func (tiny *Tiny) Run(port string) {
	log.Println("Tiny is listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, tiny))
}

// 添加中间件
func (tiny *Tiny) Use(h Handler) {
	tiny.middlewares = append(tiny.middlewares, h)
}

// 设置路由匹配失败时的Handler。如果不设置，则在匹配失败时返回404状态码。
func (tiny *Tiny) NotFound(h Handler) {
	tiny.notFoundHandler = h
}

// 设置处理panic的Handler。如果不设置，则在某个路由发生没有recover的panic时，返回500状态码。
func (tiny *Tiny) PanicHandler(h Handler) {
	tiny.panicHandler = h
}

// 添加处理该路由POST请求的handlers
func (tiny *Tiny) Post(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.post(handlers)
}

// 添加处理该路由GET请求的handlers
func (tiny *Tiny) Get(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.get(handlers)
}

// 添加处理该路由PUT请求的handlers
func (tiny *Tiny) Put(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.put(handlers)
}

// 添加处理该路由PATCH请求的handlers
func (tiny *Tiny) Patch(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.patch(handlers)
}

// 添加处理该路由DELETE请求的handlers
func (tiny *Tiny) Delete(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.delete(handlers)
}

// 添加处理该路由HEAD请求的handlers
func (tiny *Tiny) Head(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.head(handlers)
}

// 添加处理该路由OPTIONS请求的handlers
func (tiny *Tiny) Options(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.options(handlers)
}

// 添加处理该路由所有请求的handlers
func (tiny *Tiny) All(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.all(handlers)
}
