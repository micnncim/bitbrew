package cmd

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/pkg/cli/ui"
	"github.com/micnncim/bitbrew/pkg/config"
)

var (
	syncFunc = sync
)

func Sync(c *cli.Context) error {
	return syncFunc()
}

func sync() error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := initBitbrewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	installed, uninstalled, err := client.Sync()
	if err != nil {
		return err
	}

	for _, p := range installed {
		ui.Infof("✔ %s installed!\n", p.Filename)
	}
	for _, p := range uninstalled {
		ui.Warnf("✔ %s uninstalled!\n", p.Filename)
	}

	return nil
}
