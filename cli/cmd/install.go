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
	arg := c.Args().First()

	s := ui.NewSpinner("Installing...")
	s.Start()

	conf, err := config.New()
	if err != nil {
		return err
	}

	bitbrew, err := di.InitBitBrew(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	// Check whether the plugin installed
	if err := bitbrew.Load(); err != nil {
		return err
	}
	for _, p := range bitbrew.Plugins() {
		if arg == p.Filename {
			s.Stop()
			ui.Warnf("%s already installed\n", arg)
			return nil
		}
	}

	ctx := context.Background()
	plugins, err := bitbrew.SearchByFilename(ctx, arg)
	if err != nil {
		return err
	}

	s.Stop()
	if len(plugins) != 1 {
		ui.Errorf("plugin not found. need to specify accurate filename\n")
		return nil
	}

	plugin := plugins[0]
	if err := bitbrew.Install(plugin); err != nil {
		return err
	}
	ui.Infof("✔ %s installed!\n", plugin.Filename)

	return nil
}
