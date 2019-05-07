package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
)

func Search(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := bitbrew.InitClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	ctx := context.Background()
	plugins, err := client.Search(ctx, c.Args().First())
	switch err {
	case nil:
	case bitbrew.ErrPluginNotFound:
		ui.Warnf("%s\n", err)
		return nil
	default:
		return err
	}

	ui.Infof("✔ %d plugins hit\n", len(plugins))

	tableWriter := ui.NewTableWriter(os.Stdout)
	tableWriter.Show(plugins)

	return nil
}
