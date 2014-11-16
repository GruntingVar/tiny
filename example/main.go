package main

import (
	"fmt"
	"tiny"
)

func main() {
	app := tiny.New()

	app.Prepend(func(ctx *tiny.Context) {
		fmt.Println("pre1")
		ctx.Next()
	})

	app.Prepend(func(ctx *tiny.Context) {
		fmt.Println("pre2")
	})

	app.Get("/users/:id", func(ctx *tiny.Context) {
		ctx.Text(200, "your id is"+ctx.Params["id"])
		ctx.Data["test"] = "hello"
		ctx.Next()
	}, func(ctx *tiny.Context) {
		fmt.Println("I'm second handle")
		fmt.Println(ctx.Data["test"])
	})

	app.Get("/blogs/:id", func(ctx *tiny.Context) {
		ctx.Json(200, map[string]interface{}{
			"id": ctx.Params["id"],
		})
	})

	app.Get("/panic", func(ctx *tiny.Context) {
		panic("panic")
	})

	app.Append(func(ctx *tiny.Context) {
		fmt.Println("append1")
	})

	app.Append(func(ctx *tiny.Context) {
		fmt.Println("append1")
	})

	app.ErrorHandle(func(ctx *tiny.Context) {
		fmt.Println(ctx.Data["error"])
	})

	app.Run("3000")

}
