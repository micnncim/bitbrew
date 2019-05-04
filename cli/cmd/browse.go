package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func Browse(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	s := ui.NewSpinner("Searching...")
	s.Start()

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

	s.Stop()
	if len(plugins) != 1 {
		ui.Errorf("plugin not found. need to specify accurate filename")
		return nil
	}

	plugin := plugins[0]
	return open.Run(plugin.BitBarURL)
}
