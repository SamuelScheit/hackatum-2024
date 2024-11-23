package main

import (
	"bytes"
	"checkmate/database"

	"github.com/valyala/fasthttp"
)

var API_OFFERS_PATH_BUFFER = []byte("/api/offers")

// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if !bytes.Equal(ctx.Path(), API_OFFERS_PATH_BUFFER) {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	if ctx.IsGet() {
		GetHandler(ctx)
	} else if ctx.IsPost() {
		PostHandler(ctx)
	} else {
		ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
	}
}

func main() {

	database.Init()
	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}
