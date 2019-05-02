package ui

import (
	"github.com/olekukonko/tablewriter"
)

func (t *TableWriter) ExportSetTable(table *tablewriter.Table) {
	t.table = table
}
