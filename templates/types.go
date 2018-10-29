package templates

// TableTemplate defines the top level template data
type TableTemplate struct {
	Table        Table
	ReceiverName string
	PackageName  string
}

// Tables is a silce of Table
type Tables []Table

// Table defines the table
type Table struct {
	Name     string
	DBName   string
	PKExists bool
	PKName   string
	PKType   string
	Fields   []Field
	Imports  map[string]struct{}
}

// Fields is a silce of Field
type Fields []Field

// Field defines column in a table
type Field struct {
	Name       string
	Type       string
	ColumnName string
	TableName  string
	Nullable   bool
}
