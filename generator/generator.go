package generator

import (
	"bytes"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"

	"github.com/LUSHDigital/modelgen/db"
	"github.com/LUSHDigital/modelgen/sqltypes"
	"github.com/LUSHDigital/modelgen/templates"
)

const (
	tableTemplate = "x_table.go.tmpl"

	pkIdentifier = "PRI"
)

// Config represents generator configuration
type Config struct {
	Out     string
	Package string
}

// NewGenerator returns a generator
func NewGenerator(database *db.DB, cfg *Config) *Generator {
	return &Generator{
		box: templates.Box(),
		db:  database,
		cfg: cfg,
	}
}

// Generator represents the code generator
type Generator struct {
	box packr.Box
	db  *db.DB
	cfg *Config
}

// Run will generate template model files
func (g *Generator) Run() {
	tables := g.db.GetTables()
	if len(tables) == 0 {
		log.Fatal("database has no tables")
	}

	tb, err := g.box.MustBytes(tableTemplate)
	if err != nil {
		log.Fatal("cannot load table template")
	}

	t := template.Must(template.New("table").Funcs(templates.FuncMap).Parse(string(tb)))

	explained := g.db.ExplainTables(tables)
	tds := tableDefinitionsToTemplate(explained)

	os.MkdirAll(g.cfg.Out, 0777)
	generateFiles(tds, t, g.cfg.Out, g.cfg.Package)

	// copy in helpers and test suite
	copyFile(g.box, "x_helpers.go.tmpl", "x_helpers.go", "helpers", g.cfg.Out, g.cfg.Package)
	copyFile(g.box, "x_helpers_test.go.tmpl", "x_helpers_test.go", "helperstest", g.cfg.Out, g.cfg.Package)
}

func generateFiles(tables templates.Tables, t *template.Template, out, pkg string) {
	for _, table := range tables {
		tbl := templates.TableTemplate{
			Table:        table,
			ReceiverName: strings.ToLower(string(table.Name[0])),
			PackageName:  pkg,
		}
		buf := new(bytes.Buffer)
		err := t.Execute(buf, tbl)
		if err != nil {
			log.Fatal(err)
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		fpath := filepath.Join(out, table.DBName)
		writeFile(fpath, bytes.NewBuffer(formatted))
	}
}

func writeFile(path string, buf *bytes.Buffer) {
	f, err := os.Create(path + ".go")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	buf.WriteTo(f)
}

func tableDefinitionsToTemplate(explained map[string][]sqltypes.Explain) templates.Tables {
	var tbls templates.Tables
	for tableName, explain := range explained {
		t := tableDefinitionToTemplate(tableName, explain)
		tbls = append(tbls, t)
	}
	return tbls
}

func tableDefinitionToTemplate(tableName string, explains []sqltypes.Explain) templates.Table {
	table := templates.Table{
		Name:    ToPascalCase(tableName),
		DBName:  tableName,
		PKType:  "int64",
		PKName:  "id",
		Imports: make(map[string]struct{}),
	}

	for _, explain := range explains {
		f := templates.Field{
			Type:       sqltypes.AssertType(*explain.Type, *explain.Null),
			Name:       ToPascalCase(*explain.Field),
			ColumnName: strings.ToLower(*explain.Field),
			Nullable:   *explain.Null == "YES",
		}

		if *explain.Key == pkIdentifier {
			table.PKName = *explain.Field
			table.PKType = sqltypes.AssertType(*explain.Type, *explain.Null)
		}

		table.Fields = append(table.Fields, f)

		if imp, ok := sqltypes.NeedsImport(f.Type); ok {
			table.Imports[imp] = struct{}{}
		}
	}

	return table
}

func copyFile(box packr.Box, src, dst, templateName, out, pkg string) {
	dbFile, err := box.MustBytes(src)
	if err != nil {
		log.Fatalf("cannot retrieve template file: %v", err)
	}

	t := template.Must(template.New(templateName).Parse(string(dbFile)))
	buf := new(bytes.Buffer)
	err = t.Execute(buf, map[string]string{
		"PackageName": pkg,
	})
	if err != nil {
		log.Fatal(err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	buf = bytes.NewBuffer(formatted)

	if err != nil {
		log.Fatal("cannot copy file")
	}
	of := filepath.Join(out, dst)
	to, err := os.Create(of)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, buf)
	if err != nil {
		log.Fatal(err)
	}
}
