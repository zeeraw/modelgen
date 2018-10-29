package db

import (
	"database/sql"
	"log"
)

const (
	dbType = "mysql"
)

// TableDefinitions defines tables with associated comments
type TableDefinitions map[string]string

// Connect creates a new database connection
func Connect(host, dbname string) (*DB, error) {
	conn, err := sql.Open(dbType, makeDSN(host, dbname))
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		DB:   conn,
		Name: dbname,
	}, nil
}

// DB represents an sql database connection in the context of modelgen
type DB struct {
	*sql.DB
	Name string
}

// GetTables will return the table names from your database
func (db *DB) GetTables() TableDefinitions {
	tds := make(TableDefinitions)

	const stmt = `SELECT table_name, column_comment
				  FROM information_schema.columns AS c
				  WHERE c.column_key = "PRI"
				  AND c.table_schema = ?`

	rows, err := db.Query(stmt, db.Name)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		var comment string
		if err := rows.Scan(&name, &comment); err != nil {
			log.Fatal(err)
		}
		tds[name] = comment
	}
	return tds
}

func backtick(s string) string { return "`" + s + "`" }

// ExplainTable will return an explaination one database table
func (db *DB) ExplainTable(name string) []Explain {

	rows, err := db.Query("EXPLAIN " + backtick(name))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var exs []Explain
	for rows.Next() {
		var ex Explain
		if err := rows.Scan(&ex.Field, &ex.Type, &ex.Null, &ex.Key, &ex.Default, &ex.Extra); err != nil {
			log.Fatal(err)
		}
		exs = append(exs, ex)
	}

	return exs
}

// ExplainTables will return an explaination of all the database tables
func (db *DB) ExplainTables(td TableDefinitions) map[string][]Explain {
	tables := make(map[string][]Explain)
	for name := range td {
		tables[name] = db.ExplainTable(name)
	}
	return tables
}
