package database

import (
	"checkmate/optimization"
	"checkmate/types"
	"database/sql"
	_ "embed"
	"fmt"
	"strconv"
)

var queryAmounts *sql.Stmt

//go:embed sql/read/query-amounts/query-amounts.sql
var QUERY_AMOUNTS_SQL string

func initAmountQuery() {
	var err error
	queryAmounts, err = db.Prepare(QUERY_AMOUNTS_SQL)
	if err != nil {
		panic(err)
	}
}

func convertStringToUint(in string) uint {
	res, err := strconv.ParseUint(in, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(res)
}

func QueryAmount(query types.GetParams, response *types.QueryResponse) error {
	regionMin, regionMax, regionMin2, regionMax2 := optimization.GetRegionBounds(query.RegionID)

	rows, err := queryAmounts.Query(
		regionMin,
		regionMax,
		regionMin2,
		regionMax2,
		query.TimeRangeEnd,
		query.TimeRangeStart,
		query.NumberDays*1000*60*60*24,
		query.PriceRangeWidth,
		query.MinFreeKilometerWidth,
	)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	priceRanges := []types.PriceRange{}
	carTypeCounts := types.CarTypeCount{}
	seatsCounts := []types.SeatsCount{}
	freeKilometerRanges := []types.FreeKilometerRange{}
	vollkaskoCounts := types.VollkaskoCount{}

	for rows.Next() {
		var groupingType string
		var groupingValue string
		var count uint

		err := rows.Scan(&groupingType, &groupingValue, &count)
		if err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		switch groupingType {
		case "price_range":
			_uint := convertStringToUint(groupingValue)
			priceRanges = append(priceRanges, types.PriceRange{
				Start: _uint,
				End:   _uint + query.PriceRangeWidth,
				Count: count,
			})
		case "carType":
			switch groupingValue {
			case "small":
				carTypeCounts.Small = count
			case "sports":
				carTypeCounts.Sports += count
			case "luxury":
				carTypeCounts.Luxury += count
			case "family":
				carTypeCounts.Family += count
			}
		case "numberSeats":
			_uint := convertStringToUint(groupingValue)
			seatsCounts = append(seatsCounts, types.SeatsCount{
				NumberSeats: _uint,
				Count:       count,
			})
		case "freeKilometerRange":
			_uint := convertStringToUint(groupingValue)
			freeKilometerRanges = append(freeKilometerRanges, types.FreeKilometerRange{
				Start: _uint,
				End:   _uint + query.MinFreeKilometerWidth,
				Count: count,
			})
		case "hasVollkasko":
			if groupingValue == "true" {
				vollkaskoCounts.TrueCount += count
			} else {
				vollkaskoCounts.FalseCount += count
			}
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	response.PriceRanges = priceRanges
	response.CarTypeCounts = carTypeCounts
	response.FreeKilometerRange = freeKilometerRanges
	response.VollkaskoCount = vollkaskoCounts
	response.SeatsCount = seatsCounts
	return nil

}
