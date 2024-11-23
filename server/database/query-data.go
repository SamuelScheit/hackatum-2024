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

func initDataQuery() {
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
	regionMin, regionMax, regionMin2, regionMax2 := optimization.GetRegionBounds(params.RegionID)

	var query *sql.Stmt

	if params.SortOrder == types.SortOrderPriceAsc {
		query = queryAsc
	} else {
		query = queryDesc
	}

	var rows *sql.Rows
	rows, err := query.Query(
		regionMin, regionMax,
		regionMin2, regionMax2,
		params.TimeRangeEnd,
		params.TimeRangeStart,
		params.NumberDays*1000*60*60*24,
		params.MinNumberSeats,
		params.MinPrice,
		params.MaxPrice,
		params.CarType,
		params.OnlyVollkasko,
		params.MinFreeKilometer,
		params.PageSize,
		params.Page*params.PageSize,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jsonData := []byte("[")
	first := true

	for rows.Next() {
		var data []byte
		var price int

		if !first {
			jsonData = append(jsonData, commaByte...)
		}
		first = false

		err = rows.Scan(&data, &price)
		if err != nil {
			return nil, err
		}

		jsonData = append(jsonData, data...)
	}

	jsonData = append(jsonData, ']')

	return jsonData, nil

}
