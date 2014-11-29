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

type Tiny struct {
	root            *routeNode
	middlewares     []Handler // 在进行路由匹配之前执行的handle
	notFoundHandler Handler
	panicHandler    Handler
}

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

func (tiny *Tiny) Run(port string) {
	log.Println("Tiny is listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, tiny))
}

func (tiny *Tiny) Use(h Handler) {
	tiny.middlewares = append(tiny.middlewares, h)
}

func (tiny *Tiny) NotFound(h Handler) {
	tiny.notFoundHandler = h
}

func (tiny *Tiny) PanicHandler(h Handler) {
	tiny.panicHandler = h
}

func (tiny *Tiny) Post(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.post(handlers)
}

func (tiny *Tiny) Get(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.get(handlers)
}

func (tiny *Tiny) Put(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.put(handlers)
}

func (tiny *Tiny) Patch(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.patch(handlers)
}

func (tiny *Tiny) Delete(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.delete(handlers)
}

func (tiny *Tiny) Head(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.head(handlers)
}

func (tiny *Tiny) Options(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.options(handlers)
}

func (tiny *Tiny) All(url string, handlers ...Handler) {
	node := tiny.root.addUrl(url)
	node.all(handlers)
}
