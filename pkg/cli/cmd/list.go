package cmd

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/pkg/bitbrew"
	"github.com/micnncim/bitbrew/pkg/cli/ui"
	"github.com/micnncim/bitbrew/pkg/config"
)

var (
	listFunc = list
)

func List(c *cli.Context) error {
	return listFunc()
}

func list() error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := initBitbrewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}
	plugins, err := client.List()
	switch err {
	case nil:
	case bitbrew.ErrPluginNotFound:
		ui.Warnf("%s\n", err)
		return nil
	default:
		return err
	}

	for _, p := range plugins {
		ui.Printf("%s\n", p.Filename)
	}

	return nil
}
