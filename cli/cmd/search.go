package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func Search(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	s := ui.NewSpinner("Searching...")
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
	ctx := context.Background()
	plugins, err := bitbrew.Search(ctx, c.Args().First())
	if err != nil {
		return err
	}

	if len(plugins) == 0 {
		ui.Errorf("\nplugin not found")
		return nil
	}
	ui.Infof("\n✔ %d plugins hit\n", len(plugins))

	tableWriter := ui.NewTableWriter(os.Stdout)
	tableWriter.Show(plugins)

	return nil
}
