package templates

// TmplStructs is a collection on TmplStruct
type TmplStructs []TmplStruct

// TmplStruct defines the table data to pass to the models
type TmplStruct struct {
	Name      string
	TableName string
	Fields    []TmplField
	Imports   map[string]struct{}
}

// TmplField defines a table field template
type TmplField struct {
	Name       string
	Type       string
	ColumnName string
	Nullable   bool
}

// StructTmplData defines the top level struct data to pass to the models
type StructTmplData struct {
	Model       TmplStruct
	Receiver    string
	PackageName string
}

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
	Name    string
	DBName  string
	PKName  string
	PKType  string
	Fields  []Field
	Imports map[string]struct{}
}

// Field defines column in a table
type Field struct {
	Name       string
	Type       string
	ColumnName string
	Nullable   bool
}
