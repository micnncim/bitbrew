package cmd

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func List(c *cli.Context) error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	bitbrew, err := di.InitBitBrew(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}
	if err := bitbrew.Load(); err != nil {
		return err
	}

	if len(bitbrew.Plugins()) != 0 {
		ui.Errorf("no plugins\n")
		return nil
	}
	for _, p := range bitbrew.Plugins() {
		ui.Printf("%s\n", p.Filename)
	}

	return nil
}
