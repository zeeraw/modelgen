package templates

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/packr"
)

const (
	packrPath = "."
)

var FuncMap = template.FuncMap{
	"select_fields":       GetSelectFields,
	"insert_fields":       GetInsertFields,
	"insert_values":       GetInsertValues,
	"insert_args":         GetInsertArgs,
	"scan_fields":         GetScanFields,
	"update_args":         GetUpdateArgs,
	"update_values":       GetUpdateValues,
	"upsert_fields":       GetUpsertFields,
	"upsert_values":       GetUpsertValues,
	"upsert_on_duplicate": GetUpsertOnDuplicate,
	"upsert_args":         GetUpsertArgs,
	"has_int_pk":          hasIntPK,
}

// Box returns the template packer box
func Box() packr.Box {
	return packr.NewBox(packrPath)
}

func GetSelectFields(fields Fields) string {
	var parts []string
	for _, fl := range fields {
		parts = append(parts, "`"+fl.ColumnName+"`")
	}
	return strings.Join(parts, ", ")
}

func GetInsertFields(fields []Field) string {
	var parts []string
	for _, fl := range fields {
		if fl.ColumnName == "id" {
			continue
		}
		parts = append(parts, "`"+fl.ColumnName+"`")
	}
	return strings.Join(parts, ", ")
}

func GetInsertValues(fields []Field) string {
	var parts []string
	for _, fl := range fields {
		switch fl.ColumnName {
		case "id":
			continue
		case "created_at":
			parts = append(parts, "NOW()")
			continue
		default:
			parts = append(parts, "?")
		}
	}
	return strings.Join(parts, ", ")
}

func GetInsertArgs(tt TableTemplate) string {
	var parts []string
	for _, fl := range tt.Table.Fields {
		switch fl.Name {
		case "ID", "CreatedAt", tt.Table.PKName:
			continue
		}
		parts = append(parts, fmt.Sprintf("%s.%s", tt.ReceiverName, fl.Name))
	}
	if len(parts) > 0 {
		return ", " + strings.Join(parts, ", ")
	}
	return ""
}

func GetScanFields(tt TableTemplate) template.HTML {
	var parts []string
	for _, fl := range tt.Table.Fields {
		parts = append(parts, fmt.Sprintf("&%s.%s", tt.ReceiverName, fl.Name))
	}
	return template.HTML(strings.Join(parts, ", "))
}

func GetUpdateArgs(tt TableTemplate) template.HTML {
	var parts []string
	for _, fl := range tt.Table.Fields {
		switch fl.Name {
		case "ID", "CreatedAt", "UpdatedAt":
			continue
		}
		parts = append(parts, fmt.Sprintf("%s.%s", tt.ReceiverName, fl.Name))
	}
	if len(parts) > 0 {
		return template.HTML(strings.Join(parts, ", ") + ", ")
	}
	return ""
}

func GetUpdateValues(tt TableTemplate) string {
	var parts []string
	for _, fl := range tt.Table.Fields {
		switch fl.Name {
		case "ID", "CreatedAt":
			continue
		case "UpdatedAt":
			parts = append(parts, fmt.Sprintf("`%s`=UTC_TIMESTAMP()", fl.ColumnName))
		default:
			parts = append(parts, fmt.Sprintf("`%s`=?", fl.ColumnName))
		}
	}
	return strings.Join(parts, ", ")
}

func GetUpsertFields(fields []Field) string {
	var parts []string
	for _, fl := range fields {
		parts = append(parts, "`"+fl.ColumnName+"`")
	}
	return strings.Join(parts, ", ")
}

func GetUpsertValues(fields []Field) string {
	var parts []string
	for _, fl := range fields {
		switch fl.ColumnName {
		case "created_at":
			parts = append(parts, "NOW()")
			continue
		default:
			parts = append(parts, "?")
		}
	}
	return strings.Join(parts, ", ")
}

func GetUpsertOnDuplicate(tt TableTemplate) string {
	var parts []string
	for _, fl := range tt.Table.Fields {
		switch fl.Name {
		case "CreatedAt":
			continue
		case "ID":
			parts = append(parts, fmt.Sprintf("`%s`=LAST_INSERT_ID(`%s`)", fl.ColumnName, fl.ColumnName))
		case "UpdatedAt":
			parts = append(parts, fmt.Sprintf("`%s`=UTC_TIMESTAMP()", fl.ColumnName))
		default:
			parts = append(parts, fmt.Sprintf("`%s`=VALUES(`%s`)", fl.ColumnName, fl.ColumnName))
		}
	}
	return strings.Join(parts, ", ")
}

func GetUpsertArgs(tt TableTemplate) string {
	var parts []string
	for _, fl := range tt.Table.Fields {
		switch fl.Name {
		case "CreatedAt":
			continue
		}
		parts = append(parts, fmt.Sprintf("%s.%s", tt.ReceiverName, fl.Name))
	}
	return strings.Join(parts, ", ")
}

func hasIntPK(t Table) bool {
	return t.PKType == "int64"
}
