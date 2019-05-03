package cli

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/cmd"
)

const version = "0.0.0"

func New() *cli.App {
	a := cli.NewApp()
	a.Name = "bitbrew"
	a.Usage = "BitBar plugin manager"
	a.Version = version
	a.Commands = cli.Commands{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "bitbrew install FILENAME",
			Action:  cmd.Install,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "bitbrew uninstall FILENAME",
			Action:  cmd.Uninstall,
		},
		{
			Name:   "sync",
			Usage:  "bitbrew sync",
			Action: cmd.Sync,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "bitbrew list",
			Action:  cmd.List,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "bitbrew search TEXT",
			Action:  cmd.Search,
		},
		{
			Name:    "browse",
			Aliases: []string{"br"},
			Usage:   "bitbrew browse FILENAME",
			Action:  cmd.Browse,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "bitbrew config",
			Action:  cmd.Config,
		},
	}
	return a
}
