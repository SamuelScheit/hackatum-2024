package routes

import (
	"checkmate/memory"

	"github.com/valyala/fasthttp"
)

func DeleteHandler(ctx *fasthttp.RequestCtx) {

	// database.DeleteAllOffers()
	memory.DeleteAllOffers()

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(`OK`)

}
