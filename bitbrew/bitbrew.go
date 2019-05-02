package bitbrew

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

type Service interface {
	Plugins() plugin.Plugins
	Search(ctx context.Context, q string) (plugin.Plugins, error)
	SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error)
	Load() error
	Save() error
	Install(p *plugin.Plugin) error
	Uninstall(p *plugin.Plugin) error
}

type service struct {
	github       github.Service
	plugins      plugin.Plugins
	formulaPath  string
	pluginFolder string
}

var (
	ErrFormulaNotExist = errors.New("formula does not exist")
)

func NewService(gh github.Service, formulaPath, pluginFolder string) Service {
	return &service{
		github:       gh,
		formulaPath:  formulaPath,
		pluginFolder: pluginFolder,
	}
}

func (s *service) Plugins() plugin.Plugins {
	return s.plugins
}

func (s *service) Search(ctx context.Context, q string) (plugin.Plugins, error) {
	q = fmt.Sprintf("%s bitbar title desc repo:%s", q, pluginRepo)
	return s.github.Search(ctx, q)
}

func (s *service) SearchByFilename(ctx context.Context, filename string) (plugin.Plugins, error) {
	q := fmt.Sprintf("filename:%s repo:%s", filename, pluginRepo)
	return s.github.SearchByFilename(ctx, q)
}

func (s *service) Load() error {
	if !s.formulaExists() {
		return ErrFormulaNotExist
	}

	f, err := os.Open(s.formulaPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var plugins plugin.Plugins
	if err := yaml.NewDecoder(f).Decode(&plugins); err != nil {
		return err
	}
	s.plugins = plugins
	return nil
}

func (s *service) Save() error {
	if !s.formulaExists() {
		f, err := os.Create(s.formulaPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return yaml.NewEncoder(f).Encode(s.plugins)
	}

	f, err := os.OpenFile(s.formulaPath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(s.plugins)
}

func (s *service) formulaExists() bool {
	_, err := os.Stat(s.formulaPath)
	return err == nil
}

func (s *service) Install(p *plugin.Plugin) error {
	resp, err := http.Get(p.GitHubRawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(s.pluginFolder, p.Filename), buf, 0755); err != nil {
		return err
	}

	// Update formula
	if err := s.Load(); err != nil {
		if err == ErrFormulaNotExist {
			s.plugins = plugin.Plugins{p}
			return s.Save()
		}
		return err
	}
	s.plugins = append(s.plugins, p)
	return s.Save()
}

func (s *service) Uninstall(p *plugin.Plugin) error {
	if err := os.Remove(filepath.Join(s.pluginFolder, p.Filename)); err != nil {
		return err
	}

	// Update formula
	if err := s.Load(); err != nil {
		return err
	}
	// Remove uninstalled plugin
	ps := make(plugin.Plugins, 0, len(s.plugins)-1)
	for _, localPlugin := range s.plugins {
		if localPlugin.Name != p.Name {
			ps = append(ps, localPlugin)
		}
	}
	s.plugins = ps

	return s.Save()
}

func (s *service) addFormula(p *plugin.Plugin) error {
	if err := s.Load(); err != nil {
		if err == ErrFormulaNotExist {
			s.plugins = plugin.Plugins{p}
			return s.Save()
		}
		return err
	}
	s.plugins = append(s.plugins, p)
	log.Printf("%#v", s.Plugins())
	return s.Save()
}

func (s *service) removeFormula(p *plugin.Plugin) error {
	if err := s.Load(); err != nil {
		return err
	}
	ps := make(plugin.Plugins, 0, len(s.plugins)-1)
	for _, localPlugin := range s.plugins {
		if localPlugin.Filename != p.Filename {
			ps = append(ps, localPlugin)
		}
	}
	s.plugins = ps
	return s.Save()
}
