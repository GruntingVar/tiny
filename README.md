# Tiny
[![wercker status](https://app.wercker.com/status/6df44e4c942054978d3ee6998a31c8ed/m "wercker status")](https://app.wercker.com/project/bykey/6df44e4c942054978d3ee6998a31c8ed)

Tiny是一个采用Golang编写的Web开发框架，主要设计灵感来源于Express，可以轻松开发RESTful API。

## hello,world!
~~~ go
package main

import "github.com/GruntingVar/tiny"

func main() {
    app := tiny.New()

    app.Get("/hello", func(ctx *tiny.Context) {
        ctx.Text(200, "hello,world!")
    })

    app.Run("3000")
}

~~~
(框架和文档都还在紧张的建设中，敬请期待)