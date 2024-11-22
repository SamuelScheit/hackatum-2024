package main

import "github.com/valyala/fasthttp"

const (
	carTypeSmall = iota
	carTypeSports
	carTypeLuxury
	carTypeFamily
)

type GetParams struct {
	regionID             uint
	timeRangeStart       uint
	timeRangeEnd         uint
	numberDays           uint
	sortOrder            string
	page                 uint
	pageSize             uint
	priceRangeWidth      uint
	minFreeKilometerWidth uint
	minNumberSeats       {
		uint
		bool
	}
}

func GetHandler(ctx *fasthttp.RequestCtx) {
	uri := ctx.URI()
	args := uri.QueryArgs()

	// Required parameters
	regionID, errRegionID := args.GetUint("regionID")
	timeRangeStart, errTimeRangeStart := args.GetUint("timeRangeStart")
	timeRangeEnd, errTimeRangeEnd := args.GetUint("timeRangeEnd")
	numberDays, errNumberDays := args.GetUint("numberDays")
	sortOrder := string(args.Peek("sortOrder"))
	page, errPage := args.GetUint("page")
	pageSize, errPageSize := args.GetUint("pageSize")
	priceRangeWidth, errPriceRangeWidth := args.GetUint("priceRangeWidth")
	minFreeKilometerWidth, errMinFreeKilometerWidth := args.GetUint("minFreeKilometerWidth")

	// Optional parameters
	minNumberSeats := args.GetUintOrZero("minNumberSeats")
	minPrice := args.GetUintOrZero("minPrice")
	maxPrice := args.GetUintOrZero("maxPrice")
	carTypeString := string(args.Peek("carType"))
	onlyVollkasko := args.GetBool("onlyVollkasko")
	minFreeKilometer := args.GetUintOrZero("minFreeKilometer")

	// Error handling for required parameters
	if errRegionID != nil || errTimeRangeStart != nil || errTimeRangeEnd != nil ||
		errNumberDays != nil || errPage != nil || errPageSize != nil ||
		errPriceRangeWidth != nil || errMinFreeKilometerWidth != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Missing or invalid required parameters")
		return
	}

	// Validate the sortOrder parameter
	if sortOrder != "price-asc" && sortOrder != "price-desc" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Invalid value for sortOrder")
		return
	}

	carType := -1

	// Validate the carType parameter (if provided)
	switch carTypeString {
	case "small":
		carType = carTypeSmall
	case "sports":
		carType = carTypeSports
	case "luxury":
		carType = carTypeLuxury
	case "family":
		carType = carTypeFamily
	}

	_ = carType

	// Log or process the parameters (debugging purpose)
	println(
		regionID, timeRangeStart, timeRangeEnd, numberDays, sortOrder,
		page, pageSize, priceRangeWidth, minFreeKilometerWidth,
		minNumberSeats, minPrice, maxPrice, carTypeString,
		onlyVollkasko, minFreeKilometer,
	)

	// Respond or further process the data
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Parameters successfully processed")
}
