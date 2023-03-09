package gutils

import (
	"html/template"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Table struct {
	Headers    table.Row
	Rows       []table.Row
	isTerminal bool
	writer     io.Writer
}

type Options func(rTable *Table)

func newTable() *Table {
	return &Table{
		Headers: table.Row{"id", "path"},
		Rows:    []table.Row{{1, "file1"}, {2, "file2"}, {3, "file3"}},
		writer:  os.Stdout,
	}
}

func SetHeaders(headers []any) Options {
	return func(rTable *Table) {
		rTable.Headers = headers
	}
}

func SetRows(contents [][]any) Options {
	rows := make([]table.Row, 0, len(contents))

	From(func(source chan<- any) {
		for _, item := range contents {
			source <- item
		}
	}).Buffer(10).ForAll(func(pipe <-chan any) {
		item := <-pipe
		if row, ok := item.(table.Row); ok {
			rows = append(rows, row)
		}
	})

	return func(rTable *Table) {
		rTable.Rows = rows
	}
}

func SetIsTerminal() Options {
	return func(rTable *Table) {
		rTable.isTerminal = true
	}
}

func SetOutput(writer io.Writer) Options {
	return func(rTable *Table) {
		rTable.writer = writer
	}
}

func TerminalRender(options ...Options) {
	t := newTable()
	for _, item := range options {
		item(t)
	}

	tmpl := template.Must(template.
		New("").
		Funcs(map[string]interface{}{
			"table": func(tab *Table) string {
				w := table.NewWriter()
				w.AppendHeader(tab.Headers)
				w.AppendRows(tab.Rows)
				if t.isTerminal {
					return w.Render()
				}
				return w.RenderCSV()
			},
		}).
		Parse(`{{ . | table }}`))
	tbl := &Table{
		Headers: t.Headers,
		Rows:    t.Rows,
	}
	_ = tmpl.Execute(t.writer, tbl)
}
