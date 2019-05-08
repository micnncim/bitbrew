package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
)

var (
	browseFunc = browse
)

func Browse(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}
	return browseFunc(c.Args().First())
}

func browse(filename string) error {
	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := initBitbrewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	ctx := context.Background()
	switch err := client.Browse(ctx, filename); err {
	case nil:
		return nil
	case bitbrew.ErrPluginNotFound:
		ui.Errorf("%s. need to specify accurate filename\n", err)
		return nil
	default:
		return err
	}
}
