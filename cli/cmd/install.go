package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
)

var (
	installFunc = install
)

func Install(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}
	return installFunc(c.Args().First())
}

func install(filename string) error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := initBitbrewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	plugin, err := client.Install(filename)
	switch err {
	case nil:
	case bitbrew.ErrPluginNotFound:
		ui.Errorf("%s. need to specify accurate filename\n", err)
		return nil
	case bitbrew.ErrPluginInstalled:
		ui.Warnf("%s\n", err)
		return nil
	default:
		return err
	}

	ui.Infof("✔ %s installed!\n", plugin.Filename)

	return nil
}
