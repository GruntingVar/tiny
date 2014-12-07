# Tiny
[![wercker status](https://app.wercker.com/status/6df44e4c942054978d3ee6998a31c8ed/s "wercker status")](https://app.wercker.com/project/bykey/6df44e4c942054978d3ee6998a31c8ed)
[![Coverage](http://gocover.io/_badge/github.com/GruntingVar/tiny)](http://gocover.io/github.com/GruntingVar/tiny)

Tiny是一个采用Golang编写的用于构建RESTful API Server的框架，主要设计灵感来源于[Express](http://expressjs.com/)。Tiny的目标是高效、易用。

## Hello,world!
安装好[Go](http://golang.org/)并设置好[GOPATH](http://golang.org/doc/code.html#GOPATH)后，创建如下的`.go`文件。
~~~ go
package main

import "github.com/GruntingVar/tiny"

func main() {
    app := tiny.New()

    app.Get("/", func(ctx *tiny.Context) {
        ctx.Text(200, "Hello,world!")
    })

    app.Run("3000")
}
~~~

接着安装tiny：
~~~
go get github.com/GruntingVar/tiny
~~~

假设上面的`.go`文件名字叫`app.go`，那么只要输入如下命令：
~~~
go run app.go
~~~

tiny将会监听3000端口，打开浏览器访问`localhost:3000`即可在页面看到"Hello,world!"了。

## RESTful
使用Tiny可以轻松构建RESTful API Server：
~~~go
app.Post("/blogs", handler)
app.Get("/blogs/:id", handler)
app.Put("/blogs/:id", handler)
app.Patch("/blogs/:id", handler)
app.Delete("/blogs/:id", handler)
~~~

tiny支持的HTTP方法：GET、POST、PUT、PATCH、DELETE、HEAD、OPTIONS，如果需要处理其它HTTP方法发起的请求，或是希望在一个handler里处理多个方法，可以使用tiny.All方法：
~~~go
app.All("/blogs/:id", handler)
~~~

### 支持的路由类型
1. 基本类型，形如`/users/login`
2. 带参数的路由，形如`/users/:id`，可以在handler中通过ctx.Params["id"]取得类型为字符串的id值。

## Handler
Handler是形如`func(*tiny.Context)`的函数，每个路由的每个方法都可以配置多个Handler，如：
~~~go
app.Get("/blogs/:id", handler1, handler2, handler3)
~~~

### 全局的Handler：
~~~go
app.Use(handler1)
app.Use(handler2)
~~~
tiny会按顺序执行这些handler，但是请记得__在这些handler里调用*tiny.Context的Next()方法__，才会执行下一个handler:
~~~go
app.Use(func(ctx *tiny.Context) {
    ctx.Next()
})
~~~
如果handler1中没有调用Next()方法，则执行完handler1后，就不会执行handler2。Next()方法极为有用，这会在下面详细介绍。

### 错误处理：
~~~go
app.NotFound(handler) // 匹配不到相应的路由时执行此handler
app.PanicHandler(handler) // 当某个handler发生panic且并未处理时，将会执行此handler
app.ErrorHandler(handler) // 调用ctx.Error方法时会进入到此hanlder处理
~~~

## Context
tiny.Context包含Req、Res、Params、Data等属性，声明如下：
~~~go
type Context struct {
    Req    *http.Request
    Res    http.ResponseWriter
    Params tiny.matchData
    Data   map[string]interface{}
    PanicMsg error
    ErrorMsg error
    // 其它私有属性
}
~~~

### ctx.Data
在Handler中可以使用Context的Data属性实现Handler间的通信：
~~~go
app.Get("/data", func(ctx *tiny.Context) {
    ctx.Data["test"] = "hello"
    ctx.Next()
}, func(ctx *Context) {
    ctx.Text(201, ctx.Data["test"].(string)) // ctx.Data["test"] == "hello"
})
~~~

### ctx.Params
通过Params属性获取路由参数：
~~~go
app.Get("/blogs/:id", func(ctx *tiny.Context) {
    // request: GET /blogs/123
    // ctx.Params["id"] == "123"
})
~~~

### ctx.PanicMsg
在自定义的PanicHandlers中可以获取该值：
~~~go
app.PanicHandler(func(ctx *tiny.Context) {
    ctx.Text(500, ctx.PanicMsg.Error())
})
~~~

### ctx.ErrorMsg
在自定义的ErrorHandlers中可以获取该值：
~~~go
app.ErrorHandler(func(ctx *tiny.Context) {
    ctx.Text(500, ctx.ErrorMsg.Error())
})
~~~

### ctx.Text(status int, text string)
~~~go
app.Get("/text", func(ctx *tiny.Context) {
    ctx.Text(200, "hello")
    // Response:
    // Status Code: 200
    // Content-type: text/plain; charset=UTF-8
    // Body: hello
})
~~~

### ctx.Json(status int, v interface{})
~~~go
app.Get("/json", func(ctx *tiny.Context) {
    ctx.Json(200, map[string]interface{}{
        "id": 1,
        "name": "Dart",
    })
    // Response:
    // Status Code: 200
    // Content-type: application/json; charset=UTF-8
    // Body:
    // {
    //    "id": 1,
    //    "name": "Dart"
    // }
})
~~~

### ctx.Next()
在一个Handler里调用Next()方法会立即执行下一个Handler方法，在执行完毕后还会继续执行这个Handler中ctx.Next()后面的代码，这样可以充分利用go语言中的defer，轻松写出有用的路由中间件Handler。
~~~go
app.Use(func(ctx *tiny.Context) {
    // do something
    ctx.Next()
    // do another
    // ctx.Data["test"].(int) == 1
    // ctx.Data["test2"].(int) == 2
})

app.Use(func(ctx *tiny.Context) {
    ctx.Data["test"] = 1
    ctx.Next()
})

app.Use(func(ctx *tiny.Context) {
    // ctx.Data["test"].(int) == 1
    ctx.Data["test2"] = 2
})

~~~

### ctx.Error(err error)
在一个Hanlder中调用ctx.Error方法，将会进入到错误处理Handler中，如果不希望执行Error方法后面的代码，记得return:
~~~go
app.Get("/error", func(ctx *tiny.Context) {
    ctx.Error(errors.New("test error"))
    return
    // do something
})

app.ErrorHandler(func(ctx *tiny.Context) {
    ctx.Text(500, ctx.ErrorMsg.Error())
})

// GET /error, then res.Body will be "test error"
~~~


## 目标&&发展
tiny的定位是一个用于构建RESTful API Server的框架，目标是高效、易用。

为了高效，tiny的路由是由树形结构维护，将匹配的时间复杂度降至O(logn)。tiny不支持也并不打算支持正则路由，因为正则路由匹配效率不高，而且在构建RESTful API中几乎不会用到，所以不提供此功能。如果需要，可以使用一个中间件来实现。为了确保高效，tiny也没有使用reflect包，但tiny并没有丧失灵活性。

为了易用，tiny的Handler只有一个参数类型为*tiny.Context的参数，它集成了多种功能，使用起来非常方便。tiny还提供了常用的中间件，如Gzip，未来还会提供更多有用的中间件，可以开箱即用，有tiny就可以实现很多基本的功能了。现在tiny可以向响应写入文本和Json对象，但未来还会支持更多的媒体类型，例如HTML、XML、图片、文件等。

tiny目前还在开发中，还没有完成，敬请期待。

## License
待定，但一定会是开源、自由、商业友好的。
