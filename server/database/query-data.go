package database

import (
	"checkmate/optimization"
	"checkmate/types"
	"database/sql"
	_ "embed"
	"fmt"
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
	regionMin, regionMax := optimization.GetRegionBounds(params.RegionID)

	var query *sql.Stmt

	if params.SortOrder == types.SortOrderPriceAsc {
		query = queryAsc
	} else {
		query = queryDesc
	}

	fmt.Println("regionMin ", regionMin)
	fmt.Println("regionMax ", regionMax)
	fmt.Println("params", params)

	var rows *sql.Rows
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
		fmt.Println("data ", data)

		jsonData = append(jsonData, data...)
	}

	jsonData = append(jsonData, ']')

	fmt.Println("jsonData ", string(jsonData))

	return jsonData, nil

}
