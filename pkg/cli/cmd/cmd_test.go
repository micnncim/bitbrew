package cmd_test

import (
	"context"

	"github.com/micnncim/bitbrew/pkg/plugin"
)

type fakeBitbrewClient struct {
	search    func(ctx context.Context, q string) (plugin.Plugins, error)
	browse    func(ctx context.Context, filename string) error
	list      func() (plugin.Plugins, error)
	install   func(filename string) (*plugin.Plugin, error)
	uninstall func(filename string) (*plugin.Plugin, error)
	sync      func() (installed plugin.Plugins, uninstalled plugin.Plugins, err error)
}

func (c *fakeBitbrewClient) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	return c.search(ctx, q)
}

func (c *fakeBitbrewClient) Browse(ctx context.Context, filename string) error {
	return c.browse(ctx, filename)
}

func (c *fakeBitbrewClient) List() (plugin.Plugins, error) {
	return c.list()
}

func (c *fakeBitbrewClient) Install(filename string) (*plugin.Plugin, error) {
	return c.install(filename)
}

func (c *fakeBitbrewClient) Uninstall(filename string) (*plugin.Plugin, error) {
	return c.uninstall(filename)
}

func (c *fakeBitbrewClient) Sync() (installed plugin.Plugins, uninstalled plugin.Plugins, err error) {
	return c.sync()
}
