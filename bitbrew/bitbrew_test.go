package bitbrew_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/internal/testutil"
	"github.com/micnncim/bitbrew/plugin"
)

var update = flag.Bool("update", false, "update golden files")

func Test_service_Load(t *testing.T) {
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
			s := new(bitbrew.ExportService)

			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".yaml.golden")
			if *update {
				buf, err := yaml.Marshal(tc.want)
				require.NoError(t, err)
				testutil.WriteFile(t, golden, buf)
			}
			s.ExportSetFormulaPath(golden)

			err := s.Load()

			assert.Equal(t, tc.want, s.ExportPlugins())
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_service_Save(t *testing.T) {
	cases := []struct {
		name    string
		plugins plugin.Plugins
		wantErr bool
	}{
		{
			name: "save plugins",
			plugins: plugin.Plugins{
				{Name: "name", Filename: "filename"},
				{Name: "name", Filename: "filename"},
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFilePath := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".yaml")

			s := new(bitbrew.ExportService)
			s.ExportSetFormulaPath(tmpFilePath)
			s.ExportSetPlugins(tc.plugins)

			err := s.Save()
			assert.Equal(t, tc.wantErr, err != nil)

			got := testutil.ReadFile(t, tmpFilePath)

			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".yaml.golden")
			if *update {
				testutil.WriteFile(t, golden, got)
			}

			want := testutil.ReadFile(t, golden)
			assert.Equal(t, string(want), string(got))

			testutil.RemoveFile(t, tmpFilePath)
		})
	}
}

func Test_service_Install(t *testing.T) {
	cases := []struct {
		name    string
		plugin  *plugin.Plugin
		wantErr bool
	}{
		{
			name: "install script",
			plugin: &plugin.Plugin{
				Filename: "install_script.sh",
			},
			wantErr: false,
		},
	}

	tmpDir := filepath.Join("testdata", "tmp")
	defer testutil.Mkdir(t, tmpDir)()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			golden := filepath.Join("testdata", testutil.NormalizeTestName(tc.name)+".sh.golden")

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(testutil.ReadFile(t, golden))
				require.NoError(t, err)
			}))

			tc.plugin.GitHubRawURL = srv.URL
			s := new(bitbrew.ExportService)
			s.ExportSetPluginFolder(tmpDir)

			err := s.Install(tc.plugin)
			assert.Equal(t, tc.wantErr, err != nil)

			got := testutil.ReadFile(t, filepath.Join(tmpDir, tc.plugin.Filename))

			if *update {
				testutil.WriteFile(t, golden, got)
			}

			want := testutil.ReadFile(t, golden)
			assert.Equal(t, string(want), string(got))
		})
	}
}
