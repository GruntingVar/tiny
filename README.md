# Tiny [![wercker status](https://app.wercker.com/status/6df44e4c942054978d3ee6998a31c8ed/s "wercker status")](https://app.wercker.com/project/bykey/6df44e4c942054978d3ee6998a31c8ed)

Tiny是一个采用Golang编写的用于构建RESTful API的框架，主要设计灵感来源于[Express](http://expressjs.com/)。目前处于低调开发阶段，功能还在调整中，可能会有较大变化。

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
2. 带参数的路由，形如`/users/:id`，可以在handler中通过ctx.Params["id"]取得类型为字符串的id值

## Handle
Handle是形如`func(*tiny.Context)`的函数，每个路由的每个方法都可以配置多个Handle，如：
~~~go
app.Get("/blogs/:id", handler1, handler2, handler3)
~~~

### 全局的Handle：
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
app.PanicHandle(handler) // 当某个handler发生panic且并未处理时，将会执行此handler
~~~

## Context
tiny.Context包含Req、Res、Params、Data属性，声明如下：
~~~go
type Context struct {
    Req    *http.Request
    Res    http.ResponseWriter
    Params tiny.matchData
    Data   map[string]interface{}
    // 其它私有属性
}
~~~

### ctx.Data
在Handle中可以使用Context的Data属性实现Handle间的通信：
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

### ctx.Text()
~~~go
app.Get("/text", func(ctx *tiny.Context) {
    ctx.Text(200, "hello")
    // Response:
    // Status Code: 200
    // Content-type: text/plain; charset=UTF-8
    // Body: hello
})
~~~

### ctx.Json()
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
也可以关闭ident和设定字符编码：
~~~go
ctx.Json(200, map[string]interface{}{
    "id": 1,
    "name": "Dart",
}, map[string]interface{}{
    "ident": false,
    "charset": "gbk",
})
// Response:
// {"id":1,"name":"Dart"}
~~~

### ctx.Next()
在一个Handle里调用Next()方法会立即执行下一个Handle方法，在执行完毕后还会继续执行这个Handle中ctx.Next()后面的代码，这样可以充分利用go语言中的defer，轻松写出有用的路由中间件Handle。
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

## 即将到来的新特性
1. 实用的中间件，包括gzip、跨域请求处理等
2. 添加对正则路由的支持
3. 支持更多的超媒体

(敬请期待)

## License
待定
