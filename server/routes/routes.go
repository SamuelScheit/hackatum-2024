package routes

import (
	"bytes"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var API_OFFERS_PATH_BUFFER = []byte("/api/offers")

// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {

	if !bytes.Equal(ctx.Path(), API_OFFERS_PATH_BUFFER) {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	start := time.Now()

	if ctx.IsGet() {
		GetHandler(ctx)
	} else if ctx.IsPost() {
		PostHandler(ctx)
	} else if ctx.IsDelete() {
		DeleteHandler(ctx)
	} else {
		ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
	}

	fmt.Println(string(ctx.Method()), "Request took: ", time.Since(start))

	// time.Sleep(10 * time.Millisecond)
}

func Serve() {
	server := &fasthttp.Server{
		Handler:            fastHTTPHandler,
		MaxRequestBodySize: 20 * 1024 * 1024, // 20MB
	}

	err := server.ListenAndServe(":80")
	if err != nil {
		panic(err)
	}

}
