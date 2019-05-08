package cmd_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/cmd"
	"github.com/micnncim/bitbrew/plugin"
)

func Test_install(t *testing.T) {
	cases := []struct {
		name        string
		filename    string
		installFunc func(filename string) (*plugin.Plugin, error)
		wantErr     bool
	}{
		{
			name:     "install",
			filename: "filename",
			installFunc: func(filename string) (*plugin.Plugin, error) {
				return &plugin.Plugin{Filename: "filename"}, nil
			},
			wantErr: false,
		},
		{
			name:     "install func returns ErrPluginNotFound",
			filename: "filename",
			installFunc: func(filename string) (*plugin.Plugin, error) {
				return nil, bitbrew.ErrPluginNotFound
			},
			wantErr: false,
		},
		{
			name:     "install func returns error",
			filename: "filename",
			installFunc: func(filename string) (*plugin.Plugin, error) {
				return nil, errors.New("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := cmd.ExportSetInitBitbrewClient(func(token, formulaPath, pluginFolder string) (bitbrew.Client, error) {
				return &fakeBitbrewClient{
					install: tc.installFunc,
				}, nil
			})

			err := cmd.ExportInstallFunc(tc.filename)
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}
