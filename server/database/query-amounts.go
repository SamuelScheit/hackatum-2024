package database

import (
	"checkmate/optimization"
	"checkmate/types"
	"database/sql"
	_ "embed"
	"fmt"
)

var queryAmounts *sql.Stmt

//go:embed sql/read/query-amounts/query-amounts.sql
var QUERY_AMOUNTS_SQL string

func initAmountQuery() {
	queryAmounts, err := db.Prepare(QUERY_AMOUNTS_SQL)
	if err != nil {
		panic(err)
	}
	defer queryAmounts.Close()
}

func QueryAmount(query types.GetParams, response *types.QueryResponse) error {
	regionMin, regionMax := optimization.GetRegionBounds(query.RegionID)

	rows, err := queryAmounts.Query(
		regionMin,
		regionMax,
		query.TimeRangeEnd,
		query.TimeRangeStart,
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
		var fieldName string
		var value1, value2, count uint

		err := rows.Scan(&fieldName, &value1, &value2, &count)
		if err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		switch fieldName {
		case "price_range":
			priceRanges = append(priceRanges, types.PriceRange{
				Start: value1,
				End:   value2,
				Count: count,
			})
		case "carType":
			switch value1 {
			case 1:
				carTypeCounts.Small += count
			case 2:
				carTypeCounts.Sports += count
			case 3:
				carTypeCounts.Luxury += count
			case 4:
				carTypeCounts.Family += count
			}
		case "numberSeats":
			seatsCounts = append(seatsCounts, types.SeatsCount{
				NumberSeats: value1,
				Count:       count,
			})
		case "freeKilometerRange":
			freeKilometerRanges = append(freeKilometerRanges, types.FreeKilometerRange{
				Start: value1,
				End:   value2,
				Count: count,
			})
		case "hasVollkasko":
			if value1 == 1 {
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
