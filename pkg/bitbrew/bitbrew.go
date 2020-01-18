package bitbrew

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"

	"github.com/micnncim/bitbrew/pkg/github"
	"github.com/micnncim/bitbrew/pkg/plugin"
)

const (
	pluginRepo = "matryer/bitbar-plugins"
)

// Bitbrew is an interface handling BitBar plugins
type Bitbrew interface {
	Plugins() plugin.Plugins
	Search(ctx context.Context, q string) (plugin.Plugins, error)
	SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error)
	Load() error
	Save() error
	Install(p *plugin.Plugin) error
	Uninstall(p *plugin.Plugin) error
	Sync() (installed plugin.Plugins, uninstalled plugin.Plugins, err error)
}

type bitbrew struct {
	github       github.Service
	plugins      plugin.Plugins
	formulaPath  string
	pluginFolder string
}

var (
	ErrFormulaNotExist = errors.New("formula does not exist")
)

// New instantiate Bitbrew
func New(gh github.Service, formulaPath, pluginFolder string) Bitbrew {
	return &bitbrew{
		github:       gh,
		formulaPath:  formulaPath,
		pluginFolder: pluginFolder,
	}
}

// Plugins is a getter for bitbrew.plugins
func (b *bitbrew) Plugins() plugin.Plugins {
	return b.plugins
}

// Search is a wrapper for github.Search
func (b *bitbrew) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	q = fmt.Sprintf("%s bitbar title desc repo:%s", q, pluginRepo)
	return b.github.Search(ctx, q)
}

// SearchByFilename is a wrapper for github.SearchByFilename
func (b *bitbrew) SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error) {
	q := fmt.Sprintf("filename:%s repo:%s", filename, pluginRepo)
	return b.github.SearchByFilename(ctx, q)
}

// Load loads a formula file into bitbrew.plugins
func (b *bitbrew) Load() error {
	if !b.formulaExists() {
		return ErrFormulaNotExist
	}

	f, err := os.Open(b.formulaPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var plugins plugin.Plugins
	if err := yaml.NewDecoder(f).Decode(&plugins); err != nil {
		return err
	}
	b.plugins = plugins
	return nil
}

// Save saves bitbrew.plugins in a formula file
func (b *bitbrew) Save() error {
	if !b.formulaExists() {
		f, err := os.Create(b.formulaPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return yaml.NewEncoder(f).Encode(b.plugins)
	}

	f, err := os.OpenFile(b.formulaPath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(b.plugins)
}

func (b *bitbrew) formulaExists() bool {
	_, err := os.Stat(b.formulaPath)
	return err == nil
}

// Install is a wrapper for download and addFormula
func (b *bitbrew) Install(p *plugin.Plugin) error {
	if err := b.download(p); err != nil {
		return err
	}
	return b.addFormula(p)
}

// Uninstall is a wrapper for remove and removeFormula
func (b *bitbrew) Uninstall(p *plugin.Plugin) error {
	if err := b.remove(p); err != nil {
		return err
	}
	return b.removeFormula(p)
}

func (b *bitbrew) download(ps ...*plugin.Plugin) error {
	var eg errgroup.Group
	for _, p := range ps {
		p := p
		eg.Go(func() error {
			resp, err := http.Get(p.GitHubRawURL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("bad status: %s", resp.Status)
			}

			dst, err := os.Create(filepath.Join(b.pluginFolder, p.Filename))
			if err != nil {
				return err
			}

			if _, err := io.Copy(dst, resp.Body); err != nil {
				return err
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (b *bitbrew) remove(ps ...*plugin.Plugin) error {
	var eg errgroup.Group
	for _, p := range ps {
		p := p
		eg.Go(func() error {
			if err := os.Remove(filepath.Join(b.pluginFolder, p.Filename)); err != nil {
				return err
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (b *bitbrew) addFormula(p *plugin.Plugin) error {
	if err := b.Load(); err != nil {
		if err == ErrFormulaNotExist {
			b.plugins = plugin.Plugins{p}
			return b.Save()
		}
		return err
	}
	b.plugins = append(b.plugins, p)
	return b.Save()
}

func (b *bitbrew) removeFormula(p *plugin.Plugin) error {
	if err := b.Load(); err != nil {
		return err
	}
	ps := make(plugin.Plugins, 0, len(b.plugins)-1)
	for _, localPlugin := range b.plugins {
		if localPlugin.Filename != p.Filename {
			ps = append(ps, localPlugin)
		}
	}
	b.plugins = ps
	return b.Save()
}

// Sync syncs a formula file and installed plugins in local
func (b *bitbrew) Sync() (installed plugin.Plugins, uninstalled plugin.Plugins, err error) {
	var shouldInstall, shouldUninstall plugin.Plugins
	shouldInstall, shouldUninstall, err = b.diff()
	if err != nil {
		return
	}

	if err = b.download(shouldInstall...); err != nil {
		return
	}
	if err = b.remove(shouldUninstall...); err != nil {
		return
	}

	installed = shouldInstall
	uninstalled = shouldUninstall
	return
}

func (b *bitbrew) diff() (shouldInstall plugin.Plugins, shouldUninstall plugin.Plugins, err error) {
	if !b.formulaExists() {
		return
	}

	// Load plugins in formula
	var f *os.File
	f, err = os.Open(b.formulaPath)
	if err != nil {
		return
	}
	defer f.Close()
	var fp plugin.Plugins
	if err = yaml.NewDecoder(f).Decode(&fp); err != nil {
		return
	}
	formulaPlugins := make(map[string]*plugin.Plugin, len(fp))
	for _, p := range fp {
		formulaPlugins[p.Filename] = p
	}

	// Load plugins installed
	var files []os.FileInfo
	files, err = ioutil.ReadDir(b.pluginFolder)
	if err != nil {
		return
	}
	var ip plugin.Plugins
	for _, f := range files {
		ip = append(ip, &plugin.Plugin{
			Filename: f.Name(),
		})
	}
	installedPlugins := make(map[string]*plugin.Plugin, len(ip))
	for _, p := range ip {
		installedPlugins[p.Filename] = p
	}

	shouldInstall = make(plugin.Plugins, 0, len(fp))
	for filename, p := range formulaPlugins {
		if _, ok := installedPlugins[filename]; !ok {
			shouldInstall = append(shouldInstall, p)
		}
	}
	shouldUninstall = make(plugin.Plugins, 0, len(ip))
	for filename, p := range installedPlugins {
		if _, ok := formulaPlugins[filename]; !ok {
			shouldUninstall = append(shouldUninstall, p)
		}
	}

	return
}
