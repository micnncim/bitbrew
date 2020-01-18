package cmd

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/pkg/config"
)

func Config(c *cli.Context) error {
	return config.Edit()
}
