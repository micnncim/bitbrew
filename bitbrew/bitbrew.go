package bitbrew

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/micnncim/bitbrew/github"
	"github.com/micnncim/bitbrew/plugin"
)

const (
	pluginRepo = "matryer/bitbar-plugins"
)

type Bitbrew interface {
	Plugins() plugin.Plugins
	Search(ctx context.Context, q string) (plugin.Plugins, error)
	SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error)
	Load() error
	Save() error
	Install(p *plugin.Plugin) error
	Uninstall(p *plugin.Plugin) error
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

func New(gh github.Service, formulaPath, pluginFolder string) Bitbrew {
	return &bitbrew{
		github:       gh,
		formulaPath:  formulaPath,
		pluginFolder: pluginFolder,
	}
}

func (b *bitbrew) Plugins() plugin.Plugins {
	return b.plugins
}

func (b *bitbrew) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	q = fmt.Sprintf("%s bitbar title desc repo:%s", q, pluginRepo)
	return b.github.Search(ctx, q)
}

func (b *bitbrew) SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error) {
	q := fmt.Sprintf("filename:%s repo:%s", filename, pluginRepo)
	return b.github.SearchByFilename(ctx, q)
}

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

func (b *bitbrew) Install(p *plugin.Plugin) error {
	resp, err := http.Get(p.GitHubRawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(b.pluginFolder, p.Filename), buf, 0755); err != nil {
		return err
	}

	return b.addFormula(p)
}

func (b *bitbrew) Uninstall(p *plugin.Plugin) error {
	if err := os.Remove(filepath.Join(b.pluginFolder, p.Filename)); err != nil {
		return err
	}
	return b.removeFormula(p)
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
