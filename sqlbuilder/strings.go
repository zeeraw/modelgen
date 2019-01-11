package sqlbuilder

import (
	"reflect"
	"strings"
)

// BuildSelect turns a slice of strings into columns for a select statement
// Example: "`id`, `name`, `description`"
func BuildSelect(columns interface{}) string {
	columnsV := reflect.ValueOf(columns)
	var parts = make([]string, columnsV.Len())
	for i := 0; i < columnsV.Len(); i++ {
		parts[i] = "`" + columnsV.Index(i).String() + "`"
	}
	return strings.Join(parts, ", ")
}
