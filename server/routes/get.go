package routes

import (
	"checkmate/memory"
	"checkmate/types"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func GetHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()

	// PERF: if necessary, do not use allocation, instead use a params pool
	params := &types.GetParams{}
	parseError := params.ParseArgs(args)

	if parseError != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(*parseError)
		fmt.Println("error get parse", *parseError)
		return
	}

	searchResults, err := memory.QuerySearchResults(params)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal server error")
		fmt.Println("err search results", err)
		return
	}

	responseJSON, err := json.Marshal(searchResults)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal server error")
		fmt.Println("error json marshal get", err)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseJSON)

}
