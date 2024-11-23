package database

import (
	"checkmate/types"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestQueryDatabase(t *testing.T) {
	Init()

	fmt.Printf("running test\n")

	start := time.Now()

	QuerySearchResults(types.GetParams{
		RegionID:              1,
		TimeRangeStart:        0,
		TimeRangeEnd:          int64(time.Now().Unix()),
		NumberDays:            1,
		SortOrder:             types.SortOrderPriceAsc,
		Page:                  1,
		PageSize:              10,
		PriceRangeWidth:       100,
		MinFreeKilometerWidth: 100,
		MinNumberSeats:        sql.NullInt32{Int32: 1, Valid: false},
		MinPrice:              sql.NullInt32{Int32: 0, Valid: false},
		MaxPrice:              sql.NullInt32{Int32: 0, Valid: false},
		CarType:               sql.NullInt32{Int32: 0, Valid: false},
		OnlyVollkasko:         sql.NullBool{Bool: false, Valid: false},
		MinFreeKilometer:      sql.NullInt32{Int32: 0, Valid: false},
	})

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)

}
