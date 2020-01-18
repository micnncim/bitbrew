package cmd_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/pkg/bitbrew"
	"github.com/micnncim/bitbrew/pkg/cli/cmd"
	"github.com/micnncim/bitbrew/pkg/plugin"
)

func Test_list(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		listFunc func() (plugin.Plugins, error)
		wantErr  bool
	}{
		{
			name:     "list",
			filename: "filename",
			listFunc: func() (plugin.Plugins, error) {
				return plugin.Plugins{{Filename: "filename"}}, nil
			},
			wantErr: false,
		},
		{
			name:     "list func returns ErrPluginNotFound",
			filename: "filename",
			listFunc: func() (plugin.Plugins, error) {
				return nil, bitbrew.ErrPluginNotFound
			},
			wantErr: false,
		},
		{
			name:     "list func returns error",
			filename: "filename",
			listFunc: func() (plugin.Plugins, error) {
				return nil, errors.New("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := cmd.ExportSetInitBitbrewClient(func(token, formulaPath, pluginFolder string) (bitbrew.Client, error) {
				return &fakeBitbrewClient{
					list: tc.listFunc,
				}, nil
			})

			err := cmd.ExportListFunc()
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}
