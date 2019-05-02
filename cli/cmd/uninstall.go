package cmd

import (
	"fmt"

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

	s := ui.NewSpinner("Uninstalling...")
	s.Start()
	defer s.Stop()

	conf, err := config.New()
	if err != nil {
		return err
	}

	bitbrew, err := di.InitBitBrew(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	plugins, err := bitbrew.ListLocal()
	if err != nil {
		return err
	}
	for _, plugin := range plugins {
		if plugin.Filename == c.Args().First() {
			if err := bitbrew.Uninstall(plugin); err != nil {
				return err
			}
			fmt.Printf("%s uninstalled!\n", plugin)
			return nil
		}
	}

	fmt.Printf("%s not found\n", c.Args().First())
	return nil
}
