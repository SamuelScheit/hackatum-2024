package types

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

// GetParams holds the parsed query parameters
type GetParams struct {
	RegionID              uint
	TimeRangeStart        int64
	TimeRangeEnd          int64
	NumberDays            uint
	SortOrder             int
	Page                  uint
	PageSize              uint
	PriceRangeWidth       uint
	MinFreeKilometerWidth uint
	MinNumberSeats        sql.NullInt32
	MinPrice              sql.NullInt32
	MaxPrice              sql.NullInt32
	CarType               sql.NullInt32
	OnlyVollkasko         sql.NullBool
	MinFreeKilometer      sql.NullInt32
}

const (
	SortOrderPriceAsc = iota
	SortOrderPriceDesc
)

func parseUint(value []byte) (uint, error) {
	v, err := strconv.ParseUint(string(value), 10, 32)
	return uint(v), err
}
func parseint32(value []byte) (int32, error) {
	v, err := strconv.ParseInt(string(value), 10, 32)
	return int32(v), err
}

func parseint64(value []byte) (int64, error) {
	v, err := strconv.ParseInt(string(value), 10, 64)
	return v, err
}

func parseOptionalUint(value []byte) (sql.NullInt32, error) {
	if len(value) == 0 {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}, nil
	}
	v, err := parseint32(value)
	return sql.NullInt32{Int32: v, Valid: err == nil}, err
}

var carTypes = map[string]int{
	"small":  1,
	"sports": 2,
	"luxury": 3,
	"family": 4,
}

const (
	CarTypeSmall  = 1
	CarTypeSports = 2
	CarTypeLuxury = 3
	CarTypeFamily = 4
)

func (params *GetParams) ParseArgs(args *fasthttp.Args) *string {
	var parseErrors []string

	// Static byte slices for key comparisons

	args.VisitAll(func(key, value []byte) {
		fmt.Println(string(key), string(value))
	})

	params.MinNumberSeats = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	params.MinPrice = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	params.MaxPrice = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	params.CarType = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
	params.OnlyVollkasko = sql.NullBool{
		Bool:  false,
		Valid: false,
	}
	params.MinFreeKilometer = sql.NullInt32{
		Int32: 0,
		Valid: false,
	}

	handlers := map[string]func(value []byte){
		"regionID": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid regionID")
				return
			}
			params.RegionID = v
		},
		"timeRangeStart": func(value []byte) {
			v, err := parseint64(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid timeRangeStart")
				return
			}
			params.TimeRangeStart = v
		},
		"timeRangeEnd": func(value []byte) {
			v, err := parseint64(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid timeRangeEnd")
				return
			}
			params.TimeRangeEnd = v
		},
		"numberDays": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid numberDays")
				return
			}
			params.NumberDays = v
		},
		"sortOrder": func(value []byte) {
			if string(value) == "price-asc" {
				params.SortOrder = SortOrderPriceAsc
			} else if string(value) == "price-desc" {
				params.SortOrder = SortOrderPriceDesc
			} else {
				parseErrors = append(parseErrors, "Invalid sortOrder")
			}
		},
		"page": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid page")
				return
			}
			params.Page = v
		},
		"pageSize": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid pageSize")
				return
			}
			params.PageSize = v
		},
		"priceRangeWidth": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid priceRangeWidth")
				return
			}
			params.PriceRangeWidth = v
		},
		"minFreeKilometerWidth": func(value []byte) {
			v, err := parseUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minFreeKilometerWidth")
				return
			}
			params.MinFreeKilometerWidth = v
		},
		"minNumberSeats": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minNumberSeats")
				return
			}
			params.MinNumberSeats = v
		},
		"minPrice": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minPrice")
				return
			}
			params.MinPrice = v
		},
		"maxPrice": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid maxPrice")
				return
			}
			params.MaxPrice = v
		},
		"carType": func(value []byte) {
			if ct, ok := carTypes[string(value)]; ok {
				params.CarType = sql.NullInt32{
					Int32: int32(ct),
					Valid: true,
				}
			} else {
				parseErrors = append(parseErrors, "Invalid carType")
			}
		},
		"onlyVollkasko": func(value []byte) {
			params.OnlyVollkasko = sql.NullBool{
				Bool:  string(value) == "true",
				Valid: true,
			}
		},
		"minFreeKilometer": func(value []byte) {
			v, err := parseOptionalUint(value)
			if err != nil {
				parseErrors = append(parseErrors, "Invalid minFreeKilometer")
				return
			}
			params.MinFreeKilometer = v
		},
	}

	// Visit all arguments and process them
	args.VisitAll(func(key, value []byte) {
		if handler, exists := handlers[string(key)]; exists {
			handler(value)
		}
	})

	var parseError string

	for _, err := range parseErrors {
		parseError += err + ", "
	}

	return &parseError
}
