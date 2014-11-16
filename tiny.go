package tiny

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var defaultNotFoundHandle = func(ctx *Context) {
	fmt.Fprintln(ctx.Res, "Not found.")
}

var defaultErrorHandle = func(ctx *Context) {
	fmt.Fprintln(ctx.Res, "Error.")
}

type Tiny struct {
	*router
	preHandles     []Handle // 在进行路由匹配之前执行的handle
	endHandles     []Handle // 在进行路由匹配之后执行的handle
	notFoundHandle Handle
	errorHandle    Handle
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

func (tiny *Tiny) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tree, params, found := tiny.router.routeTree.find(r.URL.Path)
	ctx := &Context{r, rw, params, make(map[string]interface{}), false}

	defer func(ctx *Context) {
		if err := recover(); err != nil {
			ctx.Data["error"] = err
			tiny.errorHandle(ctx)
		}
	}(ctx)

	runHandles(ctx, tiny.preHandles)

	if found == true {
		if tree.handles["ALL"] != nil {
			runHandles(ctx, tree.handles["ALL"])
		} else {
			runHandles(ctx, tree.handles[strings.ToUpper(r.Method)])
		}
	} else {
		tiny.notFoundHandle(ctx)
	}

	runHandles(ctx, tiny.endHandles)
}

func New() *Tiny {
	return &Tiny{
		router:         newRouter(),
		preHandles:     []Handle{},
		endHandles:     []Handle{},
		notFoundHandle: defaultNotFoundHandle,
		errorHandle:    defaultErrorHandle,
	}
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
