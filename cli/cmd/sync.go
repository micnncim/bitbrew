package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/di"
)

func Sync(c *cli.Context) error {
	s := ui.NewSpinner("Syncing...")
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

	installed, uninstalled, err := bitbrew.Sync()
	if err != nil {
		return err
	}
	fmt.Println()
	for _, p := range installed {
		ui.Infof("✔ %s installed!\n", p.Filename)
	}
	for _, p := range uninstalled {
		ui.Warnf("✔ %s uninstalled!\n", p.Filename)
	}

	return nil
}
