package cmd

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
)

func List(c *cli.Context) error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := bitbrew.NewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}
	plugins, err := client.List()
	switch err {
	case nil:
	case bitbrew.ErrPluginNotFound:
		ui.Errorf("%s. need to specify accurate filename\n", err)
		return nil
	default:
		return err
	}

	for _, p := range plugins {
		ui.Printf("%s\n", p.Filename)
	}

	return nil
}
