package main

import (
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

var keys = map[string][]byte{
	"regionID":              []byte("regionID"),
	"timeRangeStart":        []byte("timeRangeStart"),
	"timeRangeEnd":          []byte("timeRangeEnd"),
	"numberDays":            []byte("numberDays"),
	"sortOrder":             []byte("sortOrder"),
	"page":                  []byte("page"),
	"pageSize":              []byte("pageSize"),
	"priceRangeWidth":       []byte("priceRangeWidth"),
	"minFreeKilometerWidth": []byte("minFreeKilometerWidth"),
	"minNumberSeats":        []byte("minNumberSeats"),
	"minPrice":              []byte("minPrice"),
	"maxPrice":              []byte("maxPrice"),
	"carType":               []byte("carType"),
	"onlyVollkasko":         []byte("onlyVollkasko"),
	"minFreeKilometer":      []byte("minFreeKilometer"),
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
			carTypeStr := string(value)
			if ct, ok := carTypes[carTypeStr]; ok {
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

func GetHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()

	// PERF: if necessary, do not use allocation, instead use a params pool
	params := GetParams{}
	parseErrors := params.parseArgs(args)

	// Handle parse errors
	if len(*parseErrors) > 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		for _, err := range *parseErrors {
			ctx.SetBodyString(err + "\n")
		}
		return
	}

	// Log or process params
	println(params.regionID, params.sortOrder, params.carType)

	// Respond
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Parameters successfully parsed")
}
