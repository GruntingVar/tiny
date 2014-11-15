package tiny

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Tiny struct {
	*router
	preHandles     []Handle // 在进行路由匹配之前执行的handle
	endHandles     []Handle // 在进行路由匹配之后执行的handle
	notFoundHandle Handle
	errorHandle    Handle
}

var defaultNotFoundHandle = func(ctx Context) {
	fmt.Fprintln(ctx.Res, "Not found.")
}

var defaultErrorHandle = func(ctx Context) {
	fmt.Fprintln(ctx.Res, "Error.")
}

func (tiny *Tiny) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tree, params, found := tiny.router.routeTree.find(r.URL.Path)
	ctx := Context{r, rw, params, make(map[string]interface{}), false}

	defer func(ctx Context) {
		if err := recover(); err != nil {
			ctx.Data["error"] = err
			tiny.errorHandle(ctx)
		}
	}(ctx)

	for _, preHandle := range tiny.preHandles {
		preHandle(ctx)
	}

	if found == true {
		if tree.handles["ALL"] != nil {
			for _, handle := range tree.handles["ALL"] {
				handle(ctx)
			}
		} else {
			for _, handle := range tree.handles[strings.ToUpper(r.Method)] {
				handle(ctx)
			}
		}
	} else {
		tiny.notFoundHandle(ctx)
	}

	for _, endHandle := range tiny.endHandles {
		endHandle(ctx)
	}
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
