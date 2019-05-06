package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
)

func Browse(c *cli.Context) error {
	if len(c.Args()) != 1 {
		return errors.New("invalid argument")
	}

	conf, err := config.New()
	if err != nil {
		return err
	}

	client, err := bitbrew.NewClient(conf.GitHub.Token, conf.BitBar.FormulaPath, conf.BitBar.PluginFolder)
	if err != nil {
		return err
	}

	ctx := context.Background()
	switch err := client.Browse(ctx, c.Args().First()); err {
	case nil:
	case bitbrew.ErrPluginNotFound:
		ui.Errorf("%s. need to specify accurate filename\n", err)
		return nil
	default:
		return err
	}

	return nil
}
