package database

import (
	"database/sql"
	_ "embed"
)

var queryAmounts *sql.Stmt

//go:embed sql/read/query-amounts/query-amounts.sql
var QUERY_AMOUNTS_SQL string
