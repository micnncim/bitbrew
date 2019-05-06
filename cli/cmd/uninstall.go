package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
)

func Uninstall(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := bitbrew.NewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	plugin, err := client.Uninstall(c.Args().First())
	switch err {
	case nil:
	case bitbrew.ErrPluginNotFound:
		ui.Errorf("%s. need to specify accurate filename\n", err)
		return nil
	default:
		return err
	}

	ui.Infof("✔ %s uninstalled!\n", plugin.Filename)

	return nil
}
