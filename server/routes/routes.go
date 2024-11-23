package routes

import (
	"bytes"
	"fmt"

	"github.com/valyala/fasthttp"
)

var API_OFFERS_PATH_BUFFER = []byte("/api/offers")

// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Println("Request URI:", string(ctx.URI().FullURI()))

	if !bytes.Equal(ctx.Path(), API_OFFERS_PATH_BUFFER) {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	if ctx.IsGet() {
		GetHandler(ctx)
	} else if ctx.IsPost() {
		PostHandler(ctx)
	} else if ctx.IsDelete() {
		DeleteHandler(ctx)
	} else {
		ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
	}
}

func Serve() {
	err := fasthttp.ListenAndServe(":8080", fastHTTPHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server started on port 8080")
}