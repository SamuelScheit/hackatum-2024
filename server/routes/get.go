package routes

import (
	"checkmate/database"
	"checkmate/types"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func GetHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()

	// PERF: if necessary, do not use allocation, instead use a params pool
	params := types.GetParams{}
	parseError := params.ParseArgs(args)

	if parseError != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(*parseError)
		fmt.Println(parseError)
		return
	}

	searchResults, err := database.QuerySearchResults(params)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal server error")
		fmt.Println(err)
		return
	}

	response := types.QueryResponse{
		Offers: types.OptimizedSearchResultOffer{
			Data: searchResults,
		},
		PriceRanges: []types.PriceRange{},
		CarTypeCounts: types.CarTypeCount{
			Small:  0,
			Sports: 0,
			Luxury: 0,
			Family: 0,
		},
		SeatsCount:         []types.SeatsCount{},
		FreeKilometerRange: []types.FreeKilometerRange{},
		VollkaskoCount: types.VollkaskoCount{
			TrueCount:  0,
			FalseCount: 0,
		},
	}

	err = database.QueryAmount(params, &response)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal server error")
		fmt.Println(err)
		return
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal server error")
		fmt.Println(err)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseJSON)

}
