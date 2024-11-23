package routes

import (
	"checkmate/database"

	"github.com/valyala/fasthttp"
)

func DeleteHandler(ctx *fasthttp.RequestCtx) {

	database.DeleteAllOffers()

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(`OK`)

}
