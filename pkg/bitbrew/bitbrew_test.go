package bitbrew_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/micnncim/bitbrew/pkg/bitbrew"
	"github.com/micnncim/bitbrew/pkg/internal/testutil"
	"github.com/micnncim/bitbrew/pkg/plugin"
)

var update = flag.Bool("update", false, "update golden files")

func Test_bitbrew_Plugins(t *testing.T) {
	cases := []struct {
		name    string
		plugins plugin.Plugins
		want    plugin.Plugins
	}{
		{
			name:    "return plugins",
			plugins: plugin.Plugins{{Name: "name"}},
			want:    plugin.Plugins{{Name: "name"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := new(bitbrew.ExportBitbrew)
			b.ExportSetPlugins(tc.plugins)

			got := b.Plugins()
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_bitbrew_Load(t *testing.T) {
	cases := []struct {
		name    string
		want    plugin.Plugins
		wantErr bool
	}{
		{
			name: "load plugins",
			want: plugin.Plugins{
				{Name: "name", Filename: "filename"},
				{Name: "name", Filename: "filename"},
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := new(bitbrew.ExportBitbrew)

			fixtureFormulaPath := filepath.Join("testdata", "fixtures", testutil.NormalizeTestName(tc.name)+".yaml")
			b.ExportSetFormulaPath(fixtureFormulaPath)

			err := b.Load()

			assert.Equal(t, tc.want, b.ExportPlugins())
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_bitbrew_Save(t *testing.T) {
	cases := []struct {
		name    string
		plugins plugin.Plugins
		wantErr bool
	}{
		{
			name: "save plugins if formula not exist",
			plugins: plugin.Plugins{
				{Name: "name", Filename: "filename"},
				{Name: "name", Filename: "filename"},
			},
			wantErr: false,
		},
		{
			name: "update formula and save plugins",
			plugins: plugin.Plugins{
				{Name: "name", Filename: "filename"},
				{Name: "name", Filename: "filename"},
				{Name: "name", Filename: "filename"},
			},
			wantErr: false,
		},
	}

	defer testutil.Mkdir(t, "tmp")()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFormulaPath := filepath.Join("tmp", testutil.NormalizeTestName(tc.name)+".yaml")

			b := new(bitbrew.ExportBitbrew)
			b.ExportSetFormulaPath(tmpFormulaPath)
			b.ExportSetPlugins(tc.plugins)

			// For updating formula case
			fixtureFormulaPath := filepath.Join("testdata", "fixtures", testutil.NormalizeTestName(tc.name)+".yaml")
			testutil.CopyFile(t, fixtureFormulaPath, tmpFormulaPath)

			err := b.Save()
			assert.Equal(t, tc.wantErr, err != nil)

			got := testutil.ReadFile(t, tmpFormulaPath)

			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".yaml.golden")
			if *update {
				testutil.WriteFile(t, golden, got)
			}

			want := testutil.ReadFile(t, golden)
			assert.Equal(t, string(want), string(got))
		})
	}
}

func Test_bitbrew_download(t *testing.T) {
	cases := []struct {
		name    string
		plugins plugin.Plugins
		wantErr bool
	}{
		{
			name: "download scripts",
			plugins: plugin.Plugins{
				{Filename: "download_scripts_1.sh"},
				{Filename: "download_scripts_2.sh"},
			},
			wantErr: false,
		},
	}

	tmpPluginFolder := "tmp"
	defer testutil.Mkdir(t, tmpPluginFolder)()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock server
			for _, p := range tc.plugins {
				golden := filepath.Join("testdata", p.Filename+".golden")

				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write(testutil.ReadFile(t, golden))
					require.NoError(t, err)
				}))
				for _, p := range tc.plugins {
					p.GitHubRawURL = srv.URL
				}
			}

			b := new(bitbrew.ExportBitbrew)
			b.ExportSetPluginFolder(tmpPluginFolder)

			err := bitbrew.ExportBitbrewDownload(b, tc.plugins...)
			assert.Equal(t, tc.wantErr, err != nil)

			for _, p := range tc.plugins {
				golden := filepath.Join("testdata", p.Filename+".golden")

				got := testutil.ReadFile(t, filepath.Join(tmpPluginFolder, p.Filename))

				if *update {
					testutil.WriteFile(t, golden, got)
				}

				want := testutil.ReadFile(t, golden)
				assert.Equal(t, string(want), string(got))
			}
		})
	}
}

func Test_bitbrew_remove(t *testing.T) {
	cases := []struct {
		name    string
		plugins plugin.Plugins
		wantErr bool
	}{
		{
			name: "remove scripts",
			plugins: plugin.Plugins{
				{Filename: "remove_scripts_1.sh"},
				{Filename: "remove_scripts_2.sh"},
			},
			wantErr: false,
		},
	}

	tmpPluginFolder := "tmp"
	defer testutil.Mkdir(t, tmpPluginFolder)()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, p := range tc.plugins {
				fixtureScript := filepath.Join("testdata", "fixtures", p.Filename)
				tmpScript := filepath.Join(tmpPluginFolder, p.Filename)
				testutil.CopyFile(t, fixtureScript, tmpScript)
			}

			b := new(bitbrew.ExportBitbrew)
			b.ExportSetPluginFolder(tmpPluginFolder)

			err := bitbrew.ExportBitbrewRemove(b, tc.plugins...)
			assert.Equal(t, tc.wantErr, err != nil)

			files := testutil.ReadDir(t, tmpPluginFolder)
			assert.Equal(t, 0, len(files))
		})
	}
}

func Test_bitbrew_addFormula(t *testing.T) {
	cases := []struct {
		name    string
		plugin  *plugin.Plugin
		wantErr bool
	}{
		{
			name:    "add plugin to formula",
			plugin:  &plugin.Plugin{Filename: "filename_added"},
			wantErr: false,
		},
		{
			name:    "create formula and add plugin",
			plugin:  &plugin.Plugin{Filename: "filename_added"},
			wantErr: false,
		},
	}

	defer testutil.Mkdir(t, "tmp")()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFormulaPath := filepath.Join("tmp", testutil.NormalizeTestName(tc.name)+".yaml")

			b := new(bitbrew.ExportBitbrew)
			b.ExportSetFormulaPath(tmpFormulaPath)

			fixtureFormulaPath := filepath.Join("testdata", "fixtures", testutil.NormalizeTestName(tc.name)+".yaml")
			testutil.CopyFile(t, fixtureFormulaPath, tmpFormulaPath)

			err := bitbrew.ExportBitbrewAddFormula(b, tc.plugin)
			assert.Equal(t, tc.wantErr, err != nil)

			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".yaml.golden")
			if *update {
				buf := testutil.ReadFile(t, tmpFormulaPath)
				testutil.WriteFile(t, golden, buf)
			}

			want := testutil.ReadFile(t, golden)
			got := testutil.ReadFile(t, tmpFormulaPath)
			assert.Equal(t, string(want), string(got))
		})
	}
}

func Test_bitbrew_removeFormula(t *testing.T) {
	cases := []struct {
		name    string
		plugin  *plugin.Plugin
		wantErr bool
	}{
		{
			name:    "remove plugin from formula",
			plugin:  &plugin.Plugin{Filename: "filename_removed"},
			wantErr: false,
		},
	}

	defer testutil.Mkdir(t, "tmp")()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFormulaPath := filepath.Join("tmp", testutil.NormalizeTestName(tc.name)+".yaml")

			b := new(bitbrew.ExportBitbrew)
			b.ExportSetFormulaPath(tmpFormulaPath)

			fixtureFormulaPath := filepath.Join("testdata", "fixtures", testutil.NormalizeTestName(tc.name)+".yaml")
			testutil.CopyFile(t, fixtureFormulaPath, tmpFormulaPath)

			err := bitbrew.ExportBitbrewRemoveFormula(b, tc.plugin)
			assert.Equal(t, tc.wantErr, err != nil)

			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".yaml.golden")
			if *update {
				buf := testutil.ReadFile(t, tmpFormulaPath)
				testutil.WriteFile(t, golden, buf)
			}

			want := testutil.ReadFile(t, golden)
			got := testutil.ReadFile(t, tmpFormulaPath)
			assert.Equal(t, string(want), string(got))
		})
	}
}

func Test_bitbrew_diff(t *testing.T) {
	cases := []struct {
		name                string
		wantShouldInstall   plugin.Plugins
		wantShouldUninstall plugin.Plugins
		wantErr             bool
	}{
		{
			name: "plus and minus",
			wantShouldInstall: plugin.Plugins{
				{Filename: "filename1.sh"},
			},
			wantShouldUninstall: plugin.Plugins{
				{Filename: "filename2.sh"},
			},
			wantErr: false,
		},
	}

	defer testutil.RemoveAll(t, filepath.Join("tmp"))

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := new(bitbrew.ExportBitbrew)

			tmpPluginFolder := filepath.Join("tmp", testutil.NormalizeTestName(tc.name))
			b.ExportSetPluginFolder(tmpPluginFolder)
			fixtureFormulaPath := filepath.Join("testdata", "fixtures", testutil.NormalizeTestName(tc.name)+".yaml")
			b.ExportSetFormulaPath(fixtureFormulaPath)

			// Copy directory
			fixturePluginFolder := filepath.Join("testdata", "fixtures", testutil.NormalizeTestName(tc.name))
			testutil.Mkdir(t, tmpPluginFolder)
			for _, p := range tc.wantShouldUninstall {
				testutil.CopyFile(t, filepath.Join(fixturePluginFolder, p.Filename), filepath.Join(tmpPluginFolder, p.Filename))
			}

			gotShouldInstall, gotShouldUninstall, err := bitbrew.ExportBitbrewDiff(b)
			assert.Equal(t, tc.wantShouldInstall, gotShouldInstall)
			assert.Equal(t, tc.wantShouldUninstall, gotShouldUninstall)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
