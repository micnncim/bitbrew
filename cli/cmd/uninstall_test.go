package cmd_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/cmd"
	"github.com/micnncim/bitbrew/plugin"
)

func Test_uninstall(t *testing.T) {
	cases := []struct {
		name          string
		filename      string
		uninstallFunc func(filename string) (*plugin.Plugin, error)
		wantErr       bool
	}{
		{
			name:     "uninstall",
			filename: "filename",
			uninstallFunc: func(filename string) (*plugin.Plugin, error) {
				return &plugin.Plugin{Filename: "filename"}, nil
			},
			wantErr: false,
		},
		{
			name:     "uninstall func returns ErrPluginNotFound",
			filename: "filename",
			uninstallFunc: func(filename string) (*plugin.Plugin, error) {
				return nil, bitbrew.ErrPluginNotFound
			},
			wantErr: false,
		},
		{
			name:     "uninstall func returns error",
			filename: "filename",
			uninstallFunc: func(filename string) (*plugin.Plugin, error) {
				return nil, errors.New("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := cmd.ExportSetInitBitbrewClient(func(token, formulaPath, pluginFolder string) (bitbrew.Client, error) {
				return &fakeBitbrewClient{
					uninstall: tc.uninstallFunc,
				}, nil
			})

			err := cmd.ExportUninstallFunc(tc.filename)
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}
