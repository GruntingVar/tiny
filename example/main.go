package main

import (
	"fmt"
	"tiny"
)

func main() {
	app := tiny.New()

	app.Prepend(func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "pre1")
	})

	app.Prepend(func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "pre2")
	})

	app.Get("/users/:id", func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "your id is", ctx.Params["id"])
		ctx.Data["test"] = "hello"
	}, func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "I'm second handle")
		fmt.Fprintln(ctx.Res, ctx.Data["test"])
	})

	app.Get("/blogs/:id", func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "blog id is", ctx.Params["id"])
	})

	app.Get("/panic", func(ctx tiny.Context) {
		panic("panic")
	})

	app.Append(func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "append1")
	})

	app.Append(func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, "append2")
	})

	app.ErrorHandle(func(ctx tiny.Context) {
		fmt.Fprintln(ctx.Res, ctx.Data["error"])
	})

	app.Run("3000")

}
