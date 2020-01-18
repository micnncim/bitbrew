package cmd_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/pkg/bitbrew"
	"github.com/micnncim/bitbrew/pkg/cli/cmd"
	"github.com/micnncim/bitbrew/pkg/plugin"
)

func Test_search(t *testing.T) {
	cases := []struct {
		name    string
		search  func(ctx context.Context, q string) (plugin.Plugins, error)
		wantErr bool
	}{
		{
			name: "search",
			search: func(ctx context.Context, q string) (plugin.Plugins, error) {
				return plugin.Plugins{{Filename: "filename"}}, nil
			},
			wantErr: false,
		},
		{
			name: "search func returns ErrPluginNotFound",
			search: func(ctx context.Context, q string) (plugin.Plugins, error) {
				return nil, bitbrew.ErrPluginNotFound
			},
			wantErr: false,
		},
		{
			name: "search func returns error",
			search: func(ctx context.Context, q string) (plugin.Plugins, error) {
				return nil, errors.New("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := cmd.ExportSetInitBitbrewClient(func(token, formulaPath, pluginFolder string) (bitbrew.Client, error) {
				return &fakeBitbrewClient{
					search: tc.search,
				}, nil
			})

			err := cmd.ExportSearchFunc("q")
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}
