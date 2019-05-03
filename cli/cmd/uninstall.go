package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func Uninstall(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	conf, err := config.New()
	if err != nil {
		return err
	}

	bitbrew, err := di.InitBitBrew(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	if err := bitbrew.Load(); err != nil {
		return err
	}
	for _, plugin := range bitbrew.Plugins() {
		if plugin.Filename == c.Args().First() {
			if err := bitbrew.Uninstall(plugin); err != nil {
				return err
			}
			ui.Printf("\n✔ %s uninstalled!\n", plugin.Filename)
			return nil
		}
	}

	ui.Errorf("\n%s not found\n", c.Args().First())
	return nil
}
