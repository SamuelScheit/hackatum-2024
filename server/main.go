package main

import (
	"bytes"

	"github.com/valyala/fasthttp"
)

var API_OFFERS_PATH_BUFFER = []byte("/api/offers")

// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if !bytes.Equal(ctx.Path(), API_OFFERS_PATH_BUFFER) {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	GetHandler(ctx)
}

func main() {

	fasthttp.ListenAndServe(":8080", fastHTTPHandler)
}
