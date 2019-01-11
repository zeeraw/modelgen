package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/LUSHDigital/modelgen/db"
	"github.com/LUSHDigital/modelgen/templates"
	"github.com/davecgh/go-spew/spew"
	"github.com/gobuffalo/packr"
	"golang.org/x/tools/imports"
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
}

func (g *Generator) copyFile(src, dst, templateName string) {
	dbFile, err := g.box.MustBytes(src)
	if err != nil {
		log.Fatalf("cannot retrieve template file: %v", err)
	}

	t := template.Must(template.New(templateName).Parse(string(dbFile)))
	buf := new(bytes.Buffer)
	err = t.Execute(buf, map[string]string{
		"PackageName": g.cfg.Package,
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
	of := filepath.Join(g.cfg.Out, dst)
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
			spew.Dump(err)
			log.Fatal(err)
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			spew.Dump(err)
			log.Fatal(err)
		}
		fmt.Println(string(formatted))

		imported, err := imports.Process("", formatted, &imports.Options{
			Comments:   true,
			AllErrors:  true,
			FormatOnly: false,
		})
		if err != nil {
			log.Fatal(err)
		}

		buf.Reset()
		buf.Write(imported)

		fpath := filepath.Join(out, table.DBName)
		writeGoFile(fpath, buf)
	}
}

func writeGoFile(path string, buf *bytes.Buffer) {
	f, err := os.Create(path + ".go")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	buf.WriteTo(f)
}

func tableDefinitionsToTemplate(explained map[string][]db.Explain) templates.Tables {
	var tbls templates.Tables
	for tableName, explain := range explained {
		t := tableDefinitionToTemplate(tableName, explain)
		tbls = append(tbls, t)
	}
	return tbls
}

func tableDefinitionToTemplate(tableName string, explains []db.Explain) templates.Table {
	table := templates.Table{
		Name:     ToPascalCase(tableName),
		DBName:   tableName,
		PKExists: false,
		PKType:   "int64",
		PKName:   "id",
		Imports:  make(map[string]struct{}),
	}

	for _, explain := range explains {
		f := templates.Field{
			Name:       ToPascalCase(*explain.Field),
			Type:       db.AssertType(*explain.Type, *explain.Null),
			ColumnName: strings.ToLower(*explain.Field),
			TableName:  table.Name,
			Nullable:   *explain.Null == "YES",
		}

		if *explain.Key == pkIdentifier {
			table.PKExists = true
			table.PKName = *explain.Field
			table.PKType = db.AssertType(*explain.Type, *explain.Null)
		}

		table.Fields = append(table.Fields, f)

		if imp, ok := db.NeedsImport(f.Type); ok {
			table.Imports[imp] = struct{}{}
		}
	}

	return table
}
