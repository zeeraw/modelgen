package modelgen

import (
	"database/sql"
)

// Querier allows sql.DB and sql.Tx to be used interchangeably, allowing you
// to use any of the model methods inside transactions or standalone calls.
type Querier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() (err error)
}

// Scanner represents a one or a collection of rows that can be scanned
type Scanner interface {
	Scan(dest ...interface{}) error
}
