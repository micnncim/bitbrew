package cmd

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/config"
)

func Config(c *cli.Context) error {
	return config.Edit()
}
