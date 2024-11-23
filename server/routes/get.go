package routes

import (
	"checkmate/database"
	"checkmate/types"
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type Response struct {
	Offers             types.OptimizedSearchResultOffer `json:"offers" validate:"required"`
	PriceRanges        []PriceRange                     `json:"priceRanges" validate:"required"`
	CarTypeCounts      CarTypeCount                     `json:"carTypeCounts" validate:"required"`
	SeatsCount         []SeatsCount                     `json:"seatsCount" validate:"required"`
	FreeKilometerRange []FreeKilometerRange             `json:"freeKilometerRange" validate:"required"`
	VollkaskoCount     VollkaskoCount                   `json:"vollkaskoCount" validate:"required"`
}

type PriceRange struct {
	Start uint `json:"start"`
	End   uint `json:"end"`
	Count uint `json:"count"`
}

type CarTypeCount struct {
	Small  uint `json:"small"`
	Sports uint `json:"sports"`
	Luxury uint `json:"luxury"`
	Family uint `json:"family"`
}

type SeatsCount struct {
	NumberSeats uint `json:"numberSeats"`
	Count       uint `json:"count"`
}

type FreeKilometerRange struct {
	Start uint `json:"start"`
	End   uint `json:"end"`
	Count uint `json:"count"`
}

type VollkaskoCount struct {
	TrueCount  uint `json:"trueCount"`
	FalseCount uint `json:"falseCount"`
}

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

	response := Response{
		Offers: types.OptimizedSearchResultOffer{
			Data: searchResults,
		},
		PriceRanges: []PriceRange{
			{Start: 1, End: 2, Count: 3},
		},
		CarTypeCounts: CarTypeCount{
			Small:  1,
			Sports: 2,
			Luxury: 3,
			Family: 4,
		},
		SeatsCount: []SeatsCount{
			{NumberSeats: 1, Count: 2},
		},
		FreeKilometerRange: []FreeKilometerRange{
			{Start: 1, End: 2, Count: 3},
		},
		VollkaskoCount: VollkaskoCount{
			TrueCount:  1,
			FalseCount: 2,
		},
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
