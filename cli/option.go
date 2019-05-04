package cli

import (
	"io"

	"github.com/urfave/cli"
)

type Option func(*cli.App)

func Writer(w io.Writer) Option {
	return func(app *cli.App) {
		app.Writer = w
	}
}

func ErrWriter(ew io.Writer) Option {
	return func(app *cli.App) {
		app.ErrWriter = ew
	}
}
