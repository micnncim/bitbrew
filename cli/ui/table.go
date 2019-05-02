package ui

import (
	"io"

	"github.com/olekukonko/tablewriter"

	"github.com/micnncim/bitbrew/plugin"
)

type TableWriter struct {
	table *tablewriter.Table
}

func NewTableWriter(writer io.Writer) *TableWriter {
	table := tablewriter.NewWriter(writer)

	table.SetHeader([]string{"Name", "Filename", "Description"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgWhiteColor},
	)

	return &TableWriter{
		table: table,
	}
}

func (t *TableWriter) Show(plugins plugin.Plugins) {
	rows := make([][]string, 0, len(plugins))
	for _, p := range plugins {
		rows = append(rows, []string{p.Name, p.Filename, p.Description})
	}
	t.table.AppendBulk(rows)
	t.table.Render()
}
