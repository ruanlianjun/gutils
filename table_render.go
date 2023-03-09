package gutils

import (
	"fmt"
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

func SetTableHeaders(headers []any) Options {
	return func(rTable *Table) {
		rTable.Headers = headers
	}
}

func SetTableRows(contents [][]any) Options {
	rows := make([]table.Row, 0, len(contents))

	for _, item := range contents {
		rows = append(rows, item)
	}

	return func(rTable *Table) {
		rTable.Rows = rows
	}
}

func SetRenderIsTerminal() Options {
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
	tbl := &Table{
		Headers: t.Headers,
		Rows:    t.Rows,
	}
	tw := table.NewWriter()
	tw.AppendHeader(tbl.Headers, table.RowConfig{
		AutoMerge: true,
	})
	tw.AppendRows(tbl.Rows)

	l := len(tbl.Rows)
	columnConfig := make([]table.ColumnConfig, 0, l)
	for i := 0; i < l; i++ {
		columnConfig = append(columnConfig,
			table.ColumnConfig{
				Number: i, Align: text.AlignLeft, AlignHeader: text.AlignCenter, AlignFooter: text.AlignCenter, WidthMin: 26, WidthMaxEnforcer: text.WrapHard,
			})
	}
	tw.SetColumnConfigs(columnConfig)
	tw.SetStyle(table.StyleLight)
	tw.Style().Options.SeparateRows = true
	tw.SetAutoIndex(true)
	fmt.Fprintf(t.writer, "%s", tw.Render())
}
