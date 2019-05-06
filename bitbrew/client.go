package bitbrew

import (
	"context"
	"errors"

	"github.com/skratchdot/open-golang/open"

	"github.com/micnncim/bitbrew/cli/ui"
	"github.com/micnncim/bitbrew/config"
	"github.com/micnncim/bitbrew/github"
	"github.com/micnncim/bitbrew/plugin"
)

type Client interface {
	Search(ctx context.Context, q string) (plugin.Plugins, error)
	Browse(ctx context.Context, filename string) error
	List() (plugin.Plugins, error)
	Install(filename string) (*plugin.Plugin, error)
	Uninstall(filename string) (*plugin.Plugin, error)
	Sync() (installed plugin.Plugins, uninstalled plugin.Plugins, err error)
}

type client struct {
	Bitbrew
}

var (
	ErrPluginNotFound  = errors.New("plugin not found")
	ErrPluginInstalled = errors.New("plugin already installed")
)

func NewClient(token, formulaPath, pluginFolder string) (Client, error) {
	if token == "" || token == config.DefaultGitHubToken {
		return nil, errors.New("github token is missing. run `bitbrew config`")
	}
	gh, err := github.NewService(token)
	if err != nil {
		return nil, err
	}
	b := New(gh, formulaPath, pluginFolder)
	return &client{Bitbrew: b}, nil
}

func (c *client) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	s := ui.NewSpinner("Searching...")
	s.Start()
	defer s.Stop()
	plugins, err := c.Bitbrew.Search(ctx, q)
	if err != nil {
		return nil, err
	}
	if len(plugins) == 0 {
		return nil, ErrPluginNotFound
	}
	return plugins, nil
}

func (c *client) Browse(ctx context.Context, filename string) error {
	s := ui.NewSpinner("Searching...")
	s.Start()
	defer s.Stop()
	plugins, err := c.Bitbrew.SearchByFilename(ctx, filename)
	if err != nil {
		return err
	}
	if len(plugins) == 0 {
		return ErrPluginNotFound
	}
	return open.Run(plugins[0].BitBarURL)
}

func (c *client) List() (plugin.Plugins, error) {
	if err := c.Bitbrew.Load(); err != nil {
		return nil, err
	}
	plugins := c.Bitbrew.Plugins()
	if len(plugins) == 0 {
		return nil, ErrPluginNotFound
	}
	return plugins, nil
}

func (c *client) Install(filename string) (*plugin.Plugin, error) {
	s := ui.NewSpinner("Installing...")
	s.Start()
	defer s.Stop()
	if err := c.Bitbrew.Load(); err != nil {
		return nil, err
	}
	for _, p := range c.Bitbrew.Plugins() {
		if p.Filename == filename {
			return nil, ErrPluginInstalled
		}
	}
	plugins, err := c.SearchByFilename(context.Background(), filename)
	if err != nil {
		return nil, err
	}
	if len(plugins) != 1 {
		return nil, ErrPluginNotFound
	}
	p := plugins[0]
	if err := c.Bitbrew.Install(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *client) Uninstall(filename string) (*plugin.Plugin, error) {
	s := ui.NewSpinner("Uninstalling...")
	s.Start()
	defer s.Stop()
	if err := c.Bitbrew.Load(); err != nil {
		return nil, err
	}
	for _, p := range c.Bitbrew.Plugins() {
		if p.Filename == filename {
			if err := c.Bitbrew.Uninstall(p); err != nil {
				return nil, err
			}
			return p, nil
		}
	}
	return nil, ErrPluginNotFound
}

func (c *client) Sync() (installed plugin.Plugins, uninstalled plugin.Plugins, err error) {
	s := ui.NewSpinner("Syncing...")
	s.Start()
	defer s.Stop()
	return c.Bitbrew.Sync()
}
