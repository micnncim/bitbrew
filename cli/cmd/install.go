package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func Install(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	s := ui.NewSpinner("Installing...")
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
	plugins, err := bitbrew.SearchByFilename(ctx, c.Args().First())
	if err != nil {
		return err
	}
	if len(plugins) != 1 {
		ui.Errorf("\nplugin not found. need to specify accurate filename\n")
		return nil
	}

	plugin := plugins[0]
	if err := bitbrew.Install(plugin); err != nil {
		return err
	}
	ui.Infof("\n✔ %s installed!\n", plugin.Filename)

	return nil
}
