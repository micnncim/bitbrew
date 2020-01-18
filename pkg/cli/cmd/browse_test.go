package cmd_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/micnncim/bitbrew/pkg/bitbrew"
	"github.com/micnncim/bitbrew/pkg/cli/cmd"
)

func Test_browse(t *testing.T) {
	cases := []struct {
		name       string
		browseFunc func(ctx context.Context, filename string) error
		wantErr    bool
	}{
		{
			name: "browse",
			browseFunc: func(ctx context.Context, filename string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "browse func returns ErrPluginNotFound",
			browseFunc: func(ctx context.Context, filename string) error {
				return bitbrew.ErrPluginNotFound
			},
			wantErr: false,
		},
		{
			name: "browse func returns error",
			browseFunc: func(ctx context.Context, filename string) error {
				return errors.New("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset := cmd.ExportSetInitBitbrewClient(func(token, formulaPath, pluginFolder string) (bitbrew.Client, error) {
				return &fakeBitbrewClient{
					browse: tc.browseFunc,
				}, nil
			})

			err := cmd.ExportBrowseFunc("filename")
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}
