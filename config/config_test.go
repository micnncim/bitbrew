package config_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/internal/testutil"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		initConfigFunc func(string) (*config.Config, error)
		want           *config.Config
		wantErr        bool
	}{
		{
			name: "new config",
			initConfigFunc: func(s string) (*config.Config, error) {
				return &config.Config{
					BitBar: &config.BitBar{
						FormulaPath: "formula_path",
					},
				}, nil
			},
			want: &config.Config{
				BitBar: &config.BitBar{
					FormulaPath: "formula_path",
				},
			},
			wantErr: false,
		},
		{
			name: "fail initConfig",
			initConfigFunc: func(s string) (*config.Config, error) {
				return nil, errors.New("error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reset := config.ExportSetInitConfigFunc(tc.initConfigFunc)

			got, err := config.New()
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}

func Test_initConfig(t *testing.T) {
	fixtures := filepath.Join("testdata", "fixtures")

	cases := []struct {
		name                string
		newDefaultViperFunc func(string) (*viper.Viper, error)
		want                *config.Config
		wantErr             bool
	}{
		{
			name: "update default config by yaml",
			newDefaultViperFunc: func(string) (*viper.Viper, error) {
				v := viper.New()
				v.SetConfigType("yaml")
				v.SetConfigName("update_default_config_by_yaml")
				v.AddConfigPath(fixtures)
				v.SetDefault("bitbar.formula_path", "formula_path")
				v.SetDefault("bitbar.plugin_folder", "plugin_folder")
				return v, nil
			},
			want: &config.Config{
				BitBar: &config.BitBar{
					FormulaPath:  "updated_formula_path",
					PluginFolder: "updated_plugin_folder",
				},
				GitHub: &config.GitHub{
					Token: "token",
				},
			},
			wantErr: false,
		},
		{
			name: "use default config",
			newDefaultViperFunc: func(string) (*viper.Viper, error) {
				v := viper.New()
				v.SetConfigType("yaml")
				v.SetConfigName("use_default_config")
				v.AddConfigPath(fixtures)
				v.SetDefault("bitbar.formula_path", "formula_path")
				v.SetDefault("bitbar.plugin_folder", "plugin_folder")
				return v, nil
			},
			want: &config.Config{
				BitBar: &config.BitBar{
					FormulaPath:  "formula_path",
					PluginFolder: "plugin_folder",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := config.ExportSetNewDefaultViperFunc(tc.newDefaultViperFunc)

			got, err := config.ExportInitConfigFunc(fixtures)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}

func Test_initConfig_CreateIfNotExist(t *testing.T) {
	// Disable Edit()
	defer config.ExportSetEditFunc(func() error { return nil })()

	tmpDir := filepath.Join("testdata", "tmp")

	cases := []struct {
		name                string
		newDefaultViperFunc func(string) (*viper.Viper, error)
		want                *config.Config
		wantErr             bool
	}{

		{
			name: "create config if not exist",
			newDefaultViperFunc: func(string) (*viper.Viper, error) {
				v := viper.New()
				v.SetConfigType("yaml")
				v.SetConfigName("create_config_if_not_exist")
				v.AddConfigPath(tmpDir)
				v.SetDefault("bitbar.formula_path", "formula_path")
				v.SetDefault("bitbar.plugin_folder", "plugin_folder")
				return v, nil
			},
			want: &config.Config{
				BitBar: &config.BitBar{
					FormulaPath:  "formula_path",
					PluginFolder: "plugin_folder",
				},
			},
			wantErr: false,
		},
	}

	defer testutil.Mkdir(t, tmpDir)()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetViper := config.ExportSetNewDefaultViperFunc(tc.newDefaultViperFunc)
			confName := testutil.NormalizeTestName(tc.name) + ".yaml"
			resetConfig := config.ExportSetDefaultConfigName(confName)

			got, err := config.ExportInitConfigFunc(tmpDir)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)

			resetViper()
			resetConfig()
		})
	}
}
