package main

import (
	"checkmate/database"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

// OptionalUint represents an optional unsigned integer
type OptionalUint struct {
	value uint
	isSet bool
}

// GetParams holds the parsed query parameters
type GetParams struct {
	regionID              uint
	timeRangeStart        uint
	timeRangeEnd          uint
	numberDays            uint
	sortOrder             string
	page                  uint
	pageSize              uint
	priceRangeWidth       uint
	minFreeKilometerWidth uint
	minNumberSeats        OptionalUint
	minPrice              OptionalUint
	maxPrice              OptionalUint
	carType               int
	onlyVollkasko         bool
	minFreeKilometer      OptionalUint
}

func parseUint(value []byte) (uint, error) {
	v, err := strconv.ParseUint(string(value), 10, 32)
	return uint(v), err
}

func parseOptionalUint(value []byte) (OptionalUint, error) {
	if len(value) == 0 {
		return OptionalUint{}, nil
	}
	v, err := parseUint(value)
	return OptionalUint{value: v, isSet: err == nil}, err
}

var carTypes = map[string]int{
	"small":  1,
	"sports": 2,
	"luxury": 3,
	"family": 4,
}

func (params *GetParams) parseArgs(args *fasthttp.Args) *[]string {
	var parseErrors []string

	// Static byte slices for key comparisons

	handlers := map[string]func(value []byte){
		"regionID": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid regionID")
				return
			}
			params.regionID = v
		},
		"timeRangeStart": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid timeRangeStart")
				return
			}
			params.timeRangeStart = v
		},
		"timeRangeEnd": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid timeRangeEnd")
				return
			}
			params.timeRangeEnd = v
		},
		"numberDays": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid numberDays")
				return
			}
			params.numberDays = v
		},
		"sortOrder": func(value []byte) {
			params.sortOrder = string(value)
		},
		"page": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid page")
				return
			}
			params.page = v
		},
		"pageSize": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid pageSize")
				return
			}
			params.pageSize = v
		},
		"priceRangeWidth": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid priceRangeWidth")
				return
			}
			params.priceRangeWidth = v
		},
		"minFreeKilometerWidth": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minFreeKilometerWidth")
				return
			}
			params.minFreeKilometerWidth = v
		},
		"minNumberSeats": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minNumberSeats")
				return
			}
			params.minNumberSeats = v
		},
		"minPrice": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minPrice")
				return
			}
			params.minPrice = v
		},
		"maxPrice": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid maxPrice")
				return
			}
			params.maxPrice = v
		},
		"carType": func(value []byte) {
			if ct, ok := carTypes[string(value)]; ok {
				params.carType = ct
			} else {
				parseErrors = append(parseErrors, "Invalid carType")
			}
		},
		"onlyVollkasko": func(value []byte) {
			params.onlyVollkasko = string(value) == "true"
		},
		"minFreeKilometer": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minFreeKilometer")
				return
			}
			params.minFreeKilometer = v
		},
	}

	// Visit all arguments and process them
	args.VisitAll(func(key, value []byte) {
		if handler, exists := handlers[string(key)]; exists {
			handler(value)
		}
	})

	return &parseErrors
}

type Response struct {
	Offers             []SearchResultOffer  `json:"offers" validate:"required"`
	PriceRanges        []PriceRange         `json:"priceRanges" validate:"required"`
	CarTypeCounts      CarTypeCount         `json:"carTypeCounts" validate:"required"`
	SeatsCount         []SeatsCount         `json:"seatsCount" validate:"required"`
	FreeKilometerRange []FreeKilometerRange `json:"freeKilometerRange" validate:"required"`
	VollkaskoCount     VollkaskoCount       `json:"vollkaskoCount" validate:"required"`
}

type SearchResultOffer struct {
	ID   string `json:"ID"`
	Data string `json:"data"`
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
	params := GetParams{}
	parseErrors := params.parseArgs(args)

	if len(*parseErrors) > 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		for _, err := range *parseErrors {
			ctx.SetBodyString(err + "\n")
		}
		fmt.Println(parseErrors)
		return
	}

	offers, err := database.RetrieveAllOffers()

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal server error")
		fmt.Println(err)
		return
	}

	searchResults := make([]SearchResultOffer, 0)

	for _, offer := range offers {
		searchResults = append(searchResults, SearchResultOffer{
			ID:   offer.ID,
			Data: offer.Data,
		})
	}

	response := Response{
		Offers: searchResults,
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
