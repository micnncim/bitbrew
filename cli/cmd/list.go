package cmd

import (
	"os"

	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func List(c *cli.Context) error {
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
	plugins, err := bitbrew.ListLocal()
	if err != nil {
		return err
	}

	tableWriter := ui.NewTableWriter(os.Stdout)
	tableWriter.Show(plugins)

	return nil
}
