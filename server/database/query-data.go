package database

import (
	"checkmate/optimization"
	"checkmate/types"
	"database/sql"
	_ "embed"
)

var queryAsc *sql.Stmt
var queryDesc *sql.Stmt

//go:embed sql/read/query-data/query-data-asc.sql
var QUERY_DATA_SQL_ASC string

//go:embed sql/read/query-data/query-data-desc.sql
var QUERY_DATA_SQL_DESC string

func initQuery() {
	var err error

	queryAsc, err = db.Prepare(QUERY_DATA_SQL_ASC)

	if err != nil {
		panic(err)
	}

	queryDesc, err = db.Prepare(QUERY_DATA_SQL_DESC)

	if err != nil {
		panic(err)
	}

}

var commaByte = []byte(",")

func QuerySearchResults(params types.GetParams) ([]byte, error) {
	regionMin, regionMax := optimization.GetRegionBounds(params.RegionID)

	var query *sql.Stmt

	if params.SortOrder == types.SortOrderPriceAsc {
		query = queryAsc
	} else {
		query = queryDesc
	}

	rows, err := query.Query(
		regionMin, regionMax,
		params.TimeRangeEnd,
		params.TimeRangeStart,
		params.NumberDays,
		0,
		params.MinNumberSeats,
		params.MinPrice,
		params.MaxPrice,
		params.CarType,
		params.OnlyVollkasko,
		params.MinFreeKilometer,
		params.PageSize,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// jsonData := []byte("[")
	first := true

	for rows.Next() {
		var data int

		if !first {
			// data = append(data, commaByte...)
		}

		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}

		// jsonData = append(jsonData, data...)
	}

	// jsonData = append(jsonData, ']')

	return nil, nil

}
