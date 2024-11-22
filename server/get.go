package checkmate

import (
	"bytes"
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

func parseArgs(args *fasthttp.Args) GetParams {

	params := &GetParams{}
	var parseErrors []string

	// Static byte slices for key comparisons

	args.VisitAll(func(key, value []byte) {
		switch {
		case bytes.Equal(key, keys["regionID"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid regionID")
				return
			}
			params.regionID = v

		case bytes.Equal(key, keys["timeRangeStart"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid timeRangeStart")
				return
			}
			params.timeRangeStart = v

		case bytes.Equal(key, keys["timeRangeEnd"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid timeRangeEnd")
				return
			}
			params.timeRangeEnd = v

		case bytes.Equal(key, keys["numberDays"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid numberDays")
				return
			}
			params.numberDays = v

		case bytes.Equal(key, keys["sortOrder"]):
			params.sortOrder = string(value)

		case bytes.Equal(key, keys["page"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid page")
				return
			}
			params.page = v

		case bytes.Equal(key, keys["pageSize"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid pageSize")
				return
			}
			params.pageSize = v

		case bytes.Equal(key, keys["priceRangeWidth"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid priceRangeWidth")
				return
			}
			params.priceRangeWidth = v

		case bytes.Equal(key, keys["minFreeKilometerWidth"]):
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minFreeKilometerWidth")
				return
			}
			params.minFreeKilometerWidth = v

		case bytes.Equal(key, keys["minNumberSeats"]):
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minNumberSeats")
				return
			}
			params.minNumberSeats = v

		case bytes.Equal(key, keys["minPrice"]):
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minPrice")
				return
			}
			params.minPrice = v

		case bytes.Equal(key, keys["maxPrice"]):
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid maxPrice")
				return
			}
			params.maxPrice = v

		case bytes.Equal(key, keys["carType"]):
			carTypeStr := string(value)
			if ct, ok := carTypes[carTypeStr]; ok {
				params.carType = ct
			} else {
				parseErrors = append(parseErrors, "Invalid carType")
			}

		case bytes.Equal(key, keys["onlyVollkasko"]):
			params.onlyVollkasko = bytes.Equal(value, []byte("true"))

		case bytes.Equal(key, keys["minFreeKilometer"]):
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minFreeKilometer")
				return
			}
			params.minFreeKilometer = v
		}
	})

	return params
}

func GetHandler(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()

	parseArgs(args)

	// Handle parse errors
	if len(parseErrors) > 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Error parsing parameters: " + strconv.Itoa(len(parseErrors)))
		return
	}

	// Log or process params
	println(params.regionID, params.sortOrder, params.carType)

	// Respond
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Parameters successfully parsed")
}
