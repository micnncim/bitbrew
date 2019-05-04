package ui_test

import (
	"bytes"
	"flag"
	"path/filepath"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/internal/testutil"
	"github.com/micnncim/bitbrew/plugin"
)

var update = flag.Bool("update", false, "update golden files")

func TestTableWriter_Show(t *testing.T) {
	cases := []struct {
		name    string
		plugins plugin.Plugins
	}{
		{
			name: "show plugins table",
			plugins: plugin.Plugins{
				{Name: "name", Filename: "filename", Description: "description"},
				{Name: "name", Filename: "filename", Description: "description"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			table := tablewriter.NewWriter(buf)
			table.SetHeader([]string{"Name", "Filename", "Description"})
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator(" ")
			table.SetColumnSeparator(" ")
			table.SetRowSeparator(" ")
			tw := new(ui.TableWriter)
			tw.ExportSetTable(table)

			tw.Show(tc.plugins)
			got := buf.Bytes()

			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".golden")
			if *update {
				testutil.WriteFile(t, golden, got)
			}

			want := testutil.ReadFile(t, golden)

			assert.Equal(t, string(want), string(got))
		})
	}
}
