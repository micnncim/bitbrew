package cli

import (
	"github.com/urfave/cli"

	"github.com/micnncim/bitbrew/cli/cmd"
)

func New() *cli.App {
	a := cli.NewApp()
	a.Name = "bitbrew"
	a.Commands = cli.Commands{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install BitBar plugin",
			Action:  cmd.Install,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "Uninstall BitBar plugin",
			Action:  cmd.Uninstall,
		},
		{
			Name:   "sync",
			Usage:  "Sync BitBar plugins with formula",
			Action: cmd.Sync,
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List installed BitBar plugins",
			Action:  cmd.List,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search BitBar plugin",
			Action:  cmd.Search,
		},
		{
			Name:    "browse",
			Aliases: []string{"br"},
			Usage:   "Browse BitBar plugins in Website",
			Action:  cmd.Browse,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Edit config file",
			Action:  cmd.Config,
		},
	}
	return a
}
