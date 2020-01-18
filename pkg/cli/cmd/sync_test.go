package cmd_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/pkg/bitbrew"
	"github.com/micnncim/bitbrew/pkg/cli/cmd"
	"github.com/micnncim/bitbrew/pkg/plugin"
)

func Test_sync(t *testing.T) {
	cases := []struct {
		name     string
		syncFunc func() (plugin.Plugins, plugin.Plugins, error)
		wantErr  bool
	}{
		{
			name: "sync",
			syncFunc: func() (plugin.Plugins, plugin.Plugins, error) {
				return plugin.Plugins{{Filename: "filename"}}, plugin.Plugins{{Filename: "filename"}}, nil
			},
			wantErr: false,
		},
		{
			name: "sync func returns error",
			syncFunc: func() (plugin.Plugins, plugin.Plugins, error) {
				return nil, nil, errors.New("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := cmd.ExportSetInitBitbrewClient(func(token, formulaPath, pluginFolder string) (bitbrew.Client, error) {
				return &fakeBitbrewClient{
					sync: tc.syncFunc,
				}, nil
			})

			err := cmd.ExportSyncFunc()
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}
