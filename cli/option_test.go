package cli_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	cliapp "github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli"
)

func TestWriter(t *testing.T) {
	cases := []struct {
		name string
		w    io.Writer
		want io.Writer
	}{
		{
			name: "set stdout",
			w:    os.Stdout,
			want: os.Stdout,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := new(cliapp.App)
			cli.Writer(tc.w)(app)
			assert.Equal(t, tc.want, app.Writer)
		})
	}
}

func TestErrWriter(t *testing.T) {
	cases := []struct {
		name string
		ew   io.Writer
		want io.Writer
	}{
		{
			name: "set stderr",
			ew:   os.Stderr,
			want: os.Stderr,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := new(cliapp.App)
			cli.ErrWriter(tc.ew)(app)
			assert.Equal(t, tc.want, app.ErrWriter)
		})
	}
}
