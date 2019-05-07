package bitbrew_test

import (
	"context"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/pkg/errors"

	"github.com/micnncim/bitbrew/bitbrew"
	"github.com/micnncim/bitbrew/plugin"
)

type fakeBitbrew struct {
	plugins          func() plugin.Plugins
	search           func(ctx context.Context, q string) (plugin.Plugins, error)
	searchByFilename func(ctx context.Context, filename string) (plugin.Plugins, error)
	load             func() error
	save             func() error
	install          func(p *plugin.Plugin) error
	uninstall        func(p *plugin.Plugin) error
	sync             func() (plugin.Plugins, plugin.Plugins, error)
}

func (b *fakeBitbrew) Plugins() plugin.Plugins {
	return b.plugins()
}

func (b *fakeBitbrew) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	return b.search(ctx, q)
}

func (b *fakeBitbrew) SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error) {
	return b.searchByFilename(ctx, filename)
}

func (b *fakeBitbrew) Load() error {
	return b.load()
}

func (b *fakeBitbrew) Save() error {
	return b.save()
}

func (b *fakeBitbrew) Install(p *plugin.Plugin) error {
	return b.install(p)
}

func (b *fakeBitbrew) Uninstall(p *plugin.Plugin) error {
	return b.uninstall(p)
}

func (b *fakeBitbrew) Sync() (plugin.Plugins, plugin.Plugins, error) {
	return b.sync()
}

func Test_client_Search(t *testing.T) {
	cases := []struct {
		name       string
		searchFunc func(ctx context.Context, q string) (plugin.Plugins, error)
		want       plugin.Plugins
		wantErr    bool
	}{
		{
			name: "found plugins",
			searchFunc: func(ctx context.Context, q string) (plugin.Plugins, error) {
				return plugin.Plugins{{Name: "name"}}, nil
			},
			want:    plugin.Plugins{{Name: "name"}},
			wantErr: false,
		},
		{
			name: "not found plugins",
			searchFunc: func(ctx context.Context, q string) (plugin.Plugins, error) {
				return nil, nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bitbrew search returns error",
			searchFunc: func(ctx context.Context, q string) (plugin.Plugins, error) {
				return nil, errors.New("error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBitbrew{
				search: tc.searchFunc,
			}
			c := bitbrew.NewClient(b)
			got, err := c.Search(context.Background(), "q")
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_client_List(t *testing.T) {
	cases := []struct {
		name        string
		loadFunc    func() error
		pluginsFunc func() plugin.Plugins
		want        plugin.Plugins
		wantErr     bool
	}{
		{
			name:     "list plugins",
			loadFunc: func() error { return nil },
			pluginsFunc: func() plugin.Plugins {
				return plugin.Plugins{{Name: "name"}}
			},
			want:    plugin.Plugins{{Name: "name"}},
			wantErr: false,
		},
		{
			name:     "load returns error",
			loadFunc: func() error { return errors.New("error") },
			want:     nil,
			wantErr:  true,
		},
		{
			name:        "plugins returns no plugins",
			loadFunc:    func() error { return nil },
			pluginsFunc: func() plugin.Plugins { return nil },
			want:        nil,
			wantErr:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBitbrew{
				plugins: tc.pluginsFunc,
				load:    tc.loadFunc,
			}
			c := bitbrew.NewClient(b)
			got, err := c.List()
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_client_Browse(t *testing.T) {
	cases := []struct {
		name                 string
		searchByFilenameFunc func(ctx context.Context, filename string) (plugin.Plugins, error)
		openFunc             func(string) error
		wantErr              bool
	}{
		{
			name: "browse",
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return plugin.Plugins{{Name: "name"}}, nil
			},
			openFunc: func(s string) error { return nil },
			wantErr:  false,
		},
		{
			name: "bitbrew search by filename returns no plugin",
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, nil
			},
			openFunc: func(s string) error { return nil },
			wantErr:  true,
		},
		{
			name: "bitbrew search by filename returns error",
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, errors.New("error")
			},
			openFunc: func(s string) error { return nil },
			wantErr:  true,
		},
		{
			name: "open func returns error",
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return plugin.Plugins{{Name: "name"}}, nil
			},
			openFunc: func(s string) error { return errors.New("error") },
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBitbrew{
				searchByFilename: tc.searchByFilenameFunc,
			}
			c := bitbrew.NewClient(b)
			reset := bitbrew.ExportSetOpenFunc(tc.openFunc)

			err := c.Browse(context.Background(), "filename")
			assert.Equal(t, tc.wantErr, err != nil)

			reset()
		})
	}
}

func Test_client_Install(t *testing.T) {
	cases := []struct {
		name                 string
		filename             string
		loadFunc             func() error
		pluginsFunc          func() plugin.Plugins
		searchByFilenameFunc func(ctx context.Context, filename string) (plugin.Plugins, error)
		installFunc          func(p *plugin.Plugin) error
		want                 *plugin.Plugin
		wantErr              bool
	}{
		{
			name:        "install plugin",
			filename:    "filename",
			loadFunc:    func() error { return nil },
			pluginsFunc: func() plugin.Plugins { return nil },
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return plugin.Plugins{{Filename: "filename"}}, nil
			},
			installFunc: func(p *plugin.Plugin) error { return nil },
			want:        &plugin.Plugin{Filename: "filename"},
			wantErr:     false,
		},
		{
			name:        "bitbrew load returns error",
			loadFunc:    func() error { return errors.New("error") },
			pluginsFunc: func() plugin.Plugins { return nil },
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, nil
			},
			installFunc: func(p *plugin.Plugin) error { return nil },
			want:        nil,
			wantErr:     true,
		},
		{
			name:     "plugin already installed",
			filename: "filename",
			loadFunc: func() error { return nil },
			pluginsFunc: func() plugin.Plugins {
				return plugin.Plugins{{Filename: "filename"}}
			},
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, nil
			},
			installFunc: func(p *plugin.Plugin) error { return nil },
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "bitbrew search by filename returns error",
			filename:    "filename",
			loadFunc:    func() error { return nil },
			pluginsFunc: func() plugin.Plugins { return nil },
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, errors.New("error")
			},
			installFunc: func(p *plugin.Plugin) error { return nil },
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "bitbrew search by filename returns no plugins",
			filename:    "filename",
			loadFunc:    func() error { return nil },
			pluginsFunc: func() plugin.Plugins { return nil },
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, nil
			},
			installFunc: func(p *plugin.Plugin) error { return nil },
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "bitbrew install returns error",
			filename:    "filename",
			loadFunc:    func() error { return nil },
			pluginsFunc: func() plugin.Plugins { return nil },
			searchByFilenameFunc: func(ctx context.Context, filename string) (plugin.Plugins, error) {
				return nil, nil
			},
			installFunc: func(p *plugin.Plugin) error { return errors.New("error") },
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBitbrew{
				load:             tc.loadFunc,
				plugins:          tc.pluginsFunc,
				searchByFilename: tc.searchByFilenameFunc,
				install:          tc.installFunc,
			}
			c := bitbrew.NewClient(b)
			got, err := c.Install(tc.filename)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_client_Uninstall(t *testing.T) {
	cases := []struct {
		name          string
		filename      string
		loadFunc      func() error
		pluginsFunc   func() plugin.Plugins
		uninstallFunc func(p *plugin.Plugin) error
		want          *plugin.Plugin
		wantErr       bool
	}{
		{
			name:     "uninstall plugin",
			filename: "filename",
			loadFunc: func() error { return nil },
			pluginsFunc: func() plugin.Plugins {
				return plugin.Plugins{{Filename: "filename"}}
			},
			uninstallFunc: func(p *plugin.Plugin) error { return nil },
			want:          &plugin.Plugin{Filename: "filename"},
			wantErr:       false,
		},
		{
			name:          "bitbrew load returns error",
			loadFunc:      func() error { return errors.New("error") },
			pluginsFunc:   func() plugin.Plugins { return nil },
			uninstallFunc: func(p *plugin.Plugin) error { return nil },
			want:          nil,
			wantErr:       true,
		},
		{
			name:     "not found plugin",
			filename: "filename",
			loadFunc: func() error { return nil },
			pluginsFunc: func() plugin.Plugins {
				return nil
			},
			uninstallFunc: func(p *plugin.Plugin) error { return nil },
			want:          nil,
			wantErr:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBitbrew{
				load:      tc.loadFunc,
				plugins:   tc.pluginsFunc,
				uninstall: tc.uninstallFunc,
			}
			c := bitbrew.NewClient(b)
			got, err := c.Uninstall(tc.filename)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func Test_client_Sync(t *testing.T) {
	cases := []struct {
		name            string
		syncFunc        func() (plugin.Plugins, plugin.Plugins, error)
		wantInstalled   plugin.Plugins
		wantUninstalled plugin.Plugins
		wantErr         bool
	}{
		{
			name: "sync",
			syncFunc: func() (plugin.Plugins, plugin.Plugins, error) {
				return plugin.Plugins{{Name: "installed"}}, plugin.Plugins{{Name: "uninstalled"}}, nil
			},
			wantInstalled:   plugin.Plugins{{Name: "installed"}},
			wantUninstalled: plugin.Plugins{{Name: "uninstalled"}},
			wantErr:         false,
		},
		{
			name: "bitbrew sync returns error",
			syncFunc: func() (plugin.Plugins, plugin.Plugins, error) {
				return nil, nil, errors.New("error")
			},
			wantInstalled:   nil,
			wantUninstalled: nil,
			wantErr:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &fakeBitbrew{
				sync: tc.syncFunc,
			}
			c := bitbrew.NewClient(b)
			gotInstalled, gotUninstalled, err := c.Sync()
			assert.Equal(t, tc.wantInstalled, gotInstalled)
			assert.Equal(t, tc.wantUninstalled, gotUninstalled)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
